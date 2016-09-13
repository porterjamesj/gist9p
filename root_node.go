package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"time"
)

type RootNode struct {
	File
	users  map[string]*UserNode
	client *github.Client
}

func NewRootNode(client *github.Client) *RootNode {
	var node RootNode
	node.File = NewDir(path(&node))
	node.client = client
	node.users = make(map[string]*UserNode)
	return &node
}

func (node *RootNode) PathComponent() string {
	return "/"
}

func (node *RootNode) Parent() FileNode {
	// the nodenode's parent is itself
	return node
}

func (node *RootNode) Child(name string) (FileNode, error) {
	// children of the root node are UserNodes
	if child, ok := node.users[name]; ok {
		return child, nil
	} else {
		user, _, err := node.client.Users.Get(name)
		if err != nil {
			return nil, err
		}
		userNode := NewUserNode(node, user)
		node.users[name] = userNode
		return userNode, nil
	}
}

func (node *RootNode) Children() ([]FileNode, error) {
	return nil, errors.New("can't list all children of root")
}

func (node *RootNode) Stat() (p9p.Dir, error) {
	now := time.Now()
	var dir = p9p.Dir{
		Mode:       0755 | p9p.DMDIR,
		AccessTime: now,
		ModTime:    now,
		Length:     0,
	}
	return dir, nil
}
