package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
)

type Node interface {
	Qid() p9p.Qid
	Open() bool
	SetOpen(bool)
	PathComponent() string
	Parent() Node
	Child(name string) (Node, error)
	Children() ([]Node, error)
	Stat() (p9p.Dir, error)
	WStat(p9p.Dir) error
	// TODO should pass the context here eventually thinking more
	// aobut it, maybe not--as far as I can think the only thing we
	// want to go with contexts is handle timeouts, so it might make
	// sense to do the waiting on the context channel in the session
	// implementation so that node implementations just don't have to
	// worry about it
	Read([]byte, int64) (int, error)
	Write([]byte, int64) (int, error)
}

type BaseNode struct {
	qid  p9p.Qid
	open bool
}

// useful default node method implementations

func (node *BaseNode) Qid() p9p.Qid {
	return node.qid
}

func (node *BaseNode) Open() bool {
	return node.open
}

func (node *BaseNode) SetOpen(value bool) {
	node.open = value
}

// to override default node method implementations

func (node *BaseNode) Child(name string) (Node, error) {
	return nil, errors.New("child not implemented")
}

func (node *BaseNode) Children() ([]Node, error) {
	return nil, errors.New("children not implemebnted")
}

func (node *BaseNode) Stat() (p9p.Dir, error) {
	return p9p.Dir{}, errors.New("stat not implemented")
}

func (node *BaseNode) WStat(p9p.Dir) error {
	return errors.New("wstat not implemented")
}

func (node *BaseNode) Read([]byte, int64) (int, error) {
	return 0, errors.New("read not implemented")
}

func (node *BaseNode) Write([]byte, int64) (int, error) {
	return 0, errors.New("write not implemented")
}

func path(node Node) string {
	parent := node.Parent()
	if parent == node {
		return node.PathComponent()
	} else {
		p := path(parent) + node.PathComponent()
		// TODO add a slash to the end if it's a directory
		return p
	}
}

const IOUNIT uint32 = 0

func open(node Node, mode p9p.Flag) error {
	// TODO verify that mode is sensible . . . this should probably
	// become a FileNode method

	// TODO I am not quite sure whether "openness" is a property of an
	// fid or the underlying file that the fid references. I suspect
	// it is the former, so I actually need to add some notion of
	// openness to the fid tracking layer (fidMap, etc.)
	node.SetOpen(true)
	return nil
}

func clunk(node Node) error {
	// TODO this treats non open fids as clunkable. i guess this is
	// fine since the man pages makes no notion of this? in practice
	// osxfuse does clunk non open fids so its sort of a moot point I
	// guess
	if node.Open() {
		node.SetOpen(false)
	}
	return nil
}

func NewFileOfType(path string, qtype p9p.QType) BaseNode {
	return BaseNode{
		qid: p9p.Qid{
			Type:    qtype,
			Version: 0,
			Path:    hashPath(path),
		},
		open: false,
	}
}

func NewDir(path string) BaseNode {
	return NewFileOfType(path, p9p.QTDIR)
}

func NewFile(path string) BaseNode {
	return NewFileOfType(path, p9p.QTFILE)
}
