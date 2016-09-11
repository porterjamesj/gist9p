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
	var root RootNode
	root.File = NewDir(path(&root))
	root.client = client
	root.children = make(map[string]*UserNode)
	return &root
}

func (root *RootNode) pathComponent() string {
	return "/"
}

func (root *RootNode) parent() FileNode {
	// the rootnode's parent is itself
	return root
}

func (root *RootNode) child(name string) (FileNode, error) {
	// children of the root node are UserNodes
	if child, ok := root.children[name]; ok {
		return child, nil
	} else {
		user, _, err := root.client.Users.Get(name)
		if err != nil {
			return nil, err
		}
		userNode := NewUserNode(root, user)
		root.children[name] = userNode
		return userNode, nil
	}
}

func (root *RootNode) stat() (p9p.Dir, error) {
	now := time.Now()
	var dir = p9p.Dir{
		Mode:       0755 | p9p.DMDIR,
		AccessTime: now,
		ModTime:    now,
		Length:     0,
	}
	return dir, nil
}
