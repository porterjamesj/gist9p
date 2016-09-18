package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"time"
)

type UserNode struct {
	BaseNode
	root   *RootNode
	user   *github.User
	gists  *GistNode
	client *github.Client
}

func NewUserNode(root *RootNode, user *github.User) *UserNode {
	var userNode UserNode
	userNode.root = root
	userNode.client = root.client
	userNode.user = user
	userNode.BaseNode = NewDir(path(&userNode))
	return &userNode
}

func (node *UserNode) Parent() Node {
	return node.root
}

func (node *UserNode) PathComponent() string {
	return *node.user.Login
}

func (node *UserNode) Child(name string) (Node, error) {
	gists, _, err := node.client.Gists.List(*node.user.Login, nil)
	if err != nil {
		return nil, err
	}
	for _, gist := range gists {
		if *gist.ID == name {
			return Node(NewGistNode(node, gist)), nil
		}
	}
	return nil, errors.New("gist not found")

}

func (node *UserNode) Children() ([]Node, error) {
	gists, _, err := node.client.Gists.List(*node.user.Login, nil)
	if err != nil {
		return nil, err
	}
	var children []Node
	for _, gist := range gists {
		// TODO these names are hellishly confusing :/
		gistNode := NewGistNode(node, gist)
		children = append(children, Node(gistNode))
	}
	return children, nil
}

func (node *UserNode) Stat() (p9p.Dir, error) {
	var times = make([]time.Time, 1)
	times[0] = node.user.CreatedAt.Time
	// TODO pagiantion, guarentee that we list all gists
	// also maybe do this by filling our list of children?
	gists, _, err := node.client.Gists.List(*node.user.Login, nil)
	if err != nil {
		return p9p.Dir{}, err
	}
	for _, gist := range gists {
		times = append(times, *gist.UpdatedAt)
	}
	modTime := maxTime(times)
	var dir = p9p.Dir{
		Mode:       0755 | p9p.DMDIR,
		AccessTime: modTime,
		ModTime:    modTime,
		// per https://swtch.com/plan9port/man/man9/stat.html,
		// "Directories and most files representing devices have a
		// conventional length of 0. "
		Length: 0,
	}
	return dir, nil
}
