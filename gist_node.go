package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
)

type GistNode struct {
	BaseNode
	user   *UserNode
	gist   *github.Gist
	client *github.Client
}

func NewGistNode(user *UserNode, gist *github.Gist) *GistNode {
	var gistNode GistNode
	gistNode.user = user
	gistNode.gist = gist
	gistNode.client = user.client
	gistNode.BaseNode = NewDir(path(&gistNode))
	return &gistNode
}

func (node *GistNode) Parent() Node {
	return node.user
}

func (node *GistNode) Child(name string) (Node, error) {
	file, ok := node.gist.Files[github.GistFilename(name)]
	if !ok {
		return nil, errors.New("gist file not found")
	}
	return NewFileNode(node, &file), nil
}

func (node *GistNode) Children() ([]Node, error) {
	var children []Node
	for _, file := range node.gist.Files {
		fileNode := NewFileNode(node, &file)
		children = append(children, Node(fileNode))
	}
	return children, nil
}

func (node *GistNode) PathComponent() string {
	return *node.gist.ID
}

func (node *GistNode) Stat() (p9p.Dir, error) {
	var dir = p9p.Dir{
		Mode:       0755 | p9p.DMDIR,
		AccessTime: *node.gist.UpdatedAt,
		ModTime:    *node.gist.UpdatedAt,
		Length:     0,
	}
	return dir, nil
}
