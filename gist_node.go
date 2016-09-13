package gist9p

import (
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
