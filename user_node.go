package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"time"
)

type UserNode struct {
	File
	root   *RootNode
	user   *github.User
	client *github.Client
}

func NewUserNode(root *RootNode, user *github.User) *UserNode {
	var userNode UserNode
	userNode.client = root.client
	userNode.user = user
	userNode.File = NewDir(path(&userNode))
	return &userNode
}

func (node *UserNode) parent() FileNode {
	return node.root
}

func (node *UserNode) pathComponent() string {
	return *node.user.Login
}

func (node *UserNode) child(name string) (FileNode, error) {
	return nil, errors.New("can't get children of users yet")
}

func (node *UserNode) stat() (p9p.Dir, error) {
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
