package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
)

type File struct {
	qid  p9p.Qid
	open bool
}

func (file *File) Qid() p9p.Qid {
	return file.qid
}

func (file *File) Open() bool {
	return file.open
}

func (file *File) SetOpen(value bool) {
	file.open = value
}

type FileNode interface {
	Qid() p9p.Qid
	Open() bool
	SetOpen(bool)
	PathComponent() string
	Parent() FileNode
	Child(name string) (FileNode, error)
	Children() ([]FileNode, error)
	Stat() (p9p.Dir, error)
}

func path(file FileNode) string {
	parent := file.Parent()
	if parent == file {
		return file.PathComponent()
	} else {
		p := path(parent) + file.PathComponent()
		// TODO add a slash to the end if it's a directory
		return p
	}
}

const IOUNIT uint32 = 0

func open(file FileNode, mode p9p.Flag) error {
	// TODO verify that mode is sensible . . . this should probably
	// become a FileNode method
	if file.Open() {
		return errors.New("file already open")
	} else {
		file.SetOpen(true)
		return nil
	}
}

func clunk(file FileNode) error {
	if file.Open() {
		file.SetOpen(false)
		return nil
	} else {
		return errors.New("file not open")
	}
}

func NewFileOfType(path string, qtype p9p.QType) File {
	return File{
		qid: p9p.Qid{
			Type:    qtype,
			Version: 0,
			Path:    hashPath(path),
		},
		open: false,
	}
}

func NewDir(path string) File {
	return NewFileOfType(path, p9p.QTDIR)
}

func NewFile(path string) File {
	return NewFileOfType(path, p9p.QTFILE)
}
