package gist9p

import (
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"log"
)

type FileNode struct {
	BaseNode
	gist    *GistNode
	file    *github.GistFile
	client  *github.Client
	content []byte
}

func NewFileNode(gist *GistNode, file *github.GistFile) *FileNode {
	var fileNode FileNode
	fileNode.gist = gist
	fileNode.file = file
	fileNode.content = nil
	fileNode.client = gist.client
	fileNode.BaseNode = NewFile(path(&fileNode))
	return &fileNode
}

func (node *FileNode) Parent() Node {
	return node.gist
}

func (node *FileNode) PathComponent() string {
	return *node.file.Filename
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
		Length:     uint64(*node.file.Size),
	}
	return dir, nil
}

func (node *FileNode) fillContent() error {
	if node.content == nil {
		err := node.gist.fillContent()
		if err == nil {
			// TODO jank, repetitive
			fname := github.GistFilename(*node.file.Filename)
			gf := node.gist.gist.Files[fname]
			node.file = &gf
			node.content = []byte(*node.file.Content)
		}
		return err
	} else {
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
	return copy(p, (*node.file.Content)[offset:]), nil
}

func (node *FileNode) Sync() error {
	// TODO this is racy. naively I would use a lock to serialize this
	// (so that editing a file / PATCHing up the result happens
	// atomically), but I guess I could learn to "share by
	// communicating" instead?
	content := string(node.content)
	node.file.Content = &content
	// TODO gross gross gross this can be gotten rid of by just making
	// this node store a reference to the github.Gist, rather than
	// copying the github.Gistfile out of it. I somehow missed the
	// fact that we were making a copy and not getting a pointer to
	// the same underlying thing
	fname := github.GistFilename(*node.file.Filename)
	node.gist.gist.Files[fname] = *node.file
	return node.gist.Sync()
}

func (node *FileNode) Write(p []byte, offset int64) (int, error) {
	var err error
	err = node.fillContent()
	if err != nil {
		log.Println("returning on error path")
		return 0, err
	}
	written := copy(node.content[offset:], p)
	if written != len(p) {
		// need to grow the internal slice and add whatever's leftover
		// TODO this feels wrong somehow, I think because I'm
		// uncertain of what it's performance characteristics will be.
		node.content = append(node.content, p[written:]...)
	}
	err = node.Sync()
	if err != nil {
		log.Println("returning on error path")
		return 0, err
	}
	log.Println("returning on success path")
	log.Println("bytes written:", written)
	return written, nil
}
