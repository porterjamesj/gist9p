package gist9p

import (
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"time"
)

type RootNode struct {
	File
	children map[string]*UserNode
	client   *github.Client
}

func NewRootNode(client *github.Client) *RootNode {
	var node RootNode
	node.File = NewDir(path(&node))
	node.client = client
	node.children = make(map[string]*UserNode)
	return &node
}

func (node *RootNode) pathComponent() string {
	return "/"
}

func (node *RootNode) parent() FileNode {
	// the nodenode's parent is itself
	return node
}

func (node *RootNode) child(name string) (FileNode, error) {
	// children of the root node are UserNodes
	if child, ok := node.children[name]; ok {
		return child, nil
	} else {
		user, _, err := node.client.Users.Get(name)
		if err != nil {
			return nil, err
		}
		userNode := NewUserNode(node, user)
		node.children[name] = userNode
		return userNode, nil
	}
}

func (node *RootNode) stat() (p9p.Dir, error) {
	now := time.Now()
	var dir = p9p.Dir{
		Mode:       0755 | p9p.DMDIR,
		AccessTime: now,
		ModTime:    now,
		Length:     0,
	}
	return dir, nil
}
