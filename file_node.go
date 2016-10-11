package gist9p

import (
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
)

type FileNode struct {
	BaseNode
	gist        *GistNode
	file        *github.GistFile
	client      *github.Client
	haveContent bool
}

func NewFileNode(gist *GistNode, file *github.GistFile) *FileNode {
	var fileNode FileNode
	fileNode.gist = gist
	fileNode.file = file
	fileNode.haveContent = false
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
	if !node.haveContent {
		err := node.gist.fillContent()
		if err == nil {
			// TODO jank, repetitive
			fname := github.GistFilename(*node.file.Filename)
			gf := node.gist.gist.Files[fname]
			node.file = &gf
			node.haveContent = true
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
