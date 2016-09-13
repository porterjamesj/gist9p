package gist9p

import (
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"time"
)

type RootNode struct {
	BaseNode
	users  map[string]*UserNode
	client *github.Client
}

func NewRootNode(client *github.Client) *RootNode {
	var node RootNode
	node.client = client
	node.users = make(map[string]*UserNode)
	node.BaseNode = NewDir(path(&node))
	return &node
}

func (node *RootNode) PathComponent() string {
	return "/"
}

func (node *RootNode) Parent() Node {
	// the nodenode's parent is itself
	return node
}

func (node *RootNode) Child(name string) (Node, error) {
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

func (node *RootNode) Children() ([]Node, error) {
	var children []Node
	for _, user := range node.users {
		children = append(children, Node(user))
	}
	return children, nil
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
