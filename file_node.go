package gist9p

import (
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"log"
	"math"
)

type FileNode struct {
	BaseNode
	gist     *GistNode
	filename github.GistFilename
	client   *github.Client
	content  []byte
}

func NewFileNode(gist *GistNode, filename github.GistFilename) *FileNode {
	var fileNode FileNode
	fileNode.gist = gist
	fileNode.filename = filename
	fileNode.content = nil
	fileNode.client = gist.client
	fileNode.BaseNode = NewFile(path(&fileNode))
	return &fileNode
}

func (node *FileNode) file() github.GistFile {
	// TODO yuck that this has to "reach up" two levels
	return node.gist.gist.Files[node.filename]
}

func (node *FileNode) Parent() Node {
	return node.gist
}

func (node *FileNode) PathComponent() string {
	return string(node.filename)
}

func (node *FileNode) WStat(dir p9p.Dir) error {
	log.Println("wstating", dir)
	log.Println("wstating length", dir.Length)
	// for now we only implement length modification
	if dir.Length != math.MaxUint64 {
		// TODO I'm not sure how to handle the case that a client is
		// trying to extend the length of a file. I don't think
		// there's a reasonable way to do this with the github API, so
		// for now we will just ignore this possibility. I also think
		// our Write implementation just sort of imagines the file is
		// always long enough so it might not matter.
		if int(dir.Length) <= len(node.content) {
			node.content = node.content[0:dir.Length]
			if dir.Length > 0 {
				// note that we don't sync if we're truncating the
				// file to have length zero, since github implicitly
				// deletes in this case. I tried using a "placeholder"
				// write here, but this led to problems since the
				// github API appears to have a bug where
				// two writes in rapid succession lead to both being
				// visible via the API, but the web interface being
				// stuck at the former.
				//
				// This truncation logic is really not sound, but so
				// far the only uses of clients truncating I've seen are
				// followed swiftly on by writes of data to replace
				// what was truncated, so this seems "fine" for now.
				//
				// One janky but possible solution would be to defer
				// syncing the content for a few seconds in another
				// goroutine, canceling if any other writes come in in
				// the meantime
				err := node.Sync()
				return err
			} else {
				return nil
			}
		} else {
			log.Println("attempt to extend file length")
			return nil
		}
	} else {
		// TODO all the non-Length stuff
		return nil
	}
}

func (node *FileNode) Stat() (p9p.Dir, error) {
	parentDir, err := node.gist.Stat()
	if err != nil {
		return p9p.Dir{}, err
	}
	var dir = p9p.Dir{
		Mode:       0755,
		AccessTime: parentDir.AccessTime,
		ModTime:    parentDir.ModTime,
		Length:     uint64(*node.file().Size),
	}
	return dir, nil
}

func (node *FileNode) fillContent() error {
	if node.content == nil {
		log.Println("about to get content from Gist ")
		err := node.gist.fillContent()
		if err == nil {
			node.content = []byte(*node.file().Content)
		}
		return err
	} else {
		log.Println("skipping doing anything")
		return nil
	}
}

func (node *FileNode) Read(p []byte, offset int64) (int, error) {
	// NOTE I'm assuming that the length of p is the amount of data
	// we're being requested to read rather than the capacity. Based
	// on what I know about go slice semantics and from reading the
	// code in 9pr, this seems correct but who knows
	//
	// TODO check for truncation?
	err := node.fillContent()
	if err != nil {
		return 0, err
	}
	return copy(p, (*node.file().Content)[offset:]), nil
}

func (node *FileNode) Sync() error {
	// TODO this is racy. naively I would use a lock to serialize this
	// (so that editing a file / PATCHing up the result happens
	// atomically), but I guess I could learn to "share by
	// communicating" instead?
	content := string(node.content)
	file := node.file()
	file.Content = &content
	node.gist.gist.Files[node.filename] = file
	return node.gist.Sync()
}

func (node *FileNode) Write(p []byte, offset int64) (int, error) {
	var err error
	log.Println("write start:", string(node.content))
	err = node.fillContent()
	log.Println("after fillcontent:", string(node.content))
	if err != nil {
		log.Println("returning on error path")
		return 0, err
	}
	neededCapacity := int(offset) + len(p)
	if len(node.content) < neededCapacity {
		// extend it
		extendBy := neededCapacity - len(node.content)
		node.content = append(node.content, make([]byte, extendBy)...)
	}
	log.Println("after extend:", string(node.content))
	written := copy(node.content[offset:], p)
	log.Println("after copy:", string(node.content))
	err = node.Sync()
	log.Println("after sync:", string(node.content))
	if err != nil {
		log.Println("returning on error path")
		return 0, err
	}
	log.Println("returning on success path")
	log.Println("bytes written:", written)
	return written, nil
}
