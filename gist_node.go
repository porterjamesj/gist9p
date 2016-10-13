package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"log"
)

type GistNode struct {
	BaseNode
	user        *UserNode
	gist        *github.Gist
	client      *github.Client
	haveContent bool
}

func NewGistNode(user *UserNode, gist *github.Gist) *GistNode {
	var gistNode GistNode
	gistNode.user = user
	gistNode.client = user.client
	gistNode.gist = gist
	gistNode.BaseNode = NewDir(path(&gistNode))
	gistNode.haveContent = false
	return &gistNode
}

func (node *GistNode) fillContent() error {
	if !node.haveContent {
		var err error
		node.gist, _, err = node.client.Gists.Get(*node.gist.ID)
		if err == nil {
			node.haveContent = true
		}
		return err
	} else {
		return nil
	}
}

func (node *GistNode) Sync() error {

	gist, resp, err := node.client.Gists.Edit(*node.gist.ID, node.gist)
	log.Println("sync PATCH resp code: ", resp.Status)
	if err != nil {
		return err
	}
	node.gist = gist
	return nil
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
