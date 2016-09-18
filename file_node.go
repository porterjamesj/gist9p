package gist9p

import (
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
)

type FileNode struct {
	BaseNode
	gist   *GistNode
	file   *github.GistFile
	client *github.Client
}

func NewFileNode(gist *GistNode, file *github.GistFile) *FileNode {
	var fileNode FileNode
	fileNode.gist = gist
	fileNode.file = file
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
	dir, nil := node.gist.Stat()
	return dir, nil
}
