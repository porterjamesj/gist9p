package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
)

type File struct {
	qid  p9p.Qid
	open bool
}

func (file *File) getQid() p9p.Qid {
	return file.qid
}

func (file *File) isOpen() bool {
	return file.open
}

func (file *File) setOpen(value bool) {
	file.open = value
}

type FileNode interface {
	getQid() p9p.Qid
	isOpen() bool
	setOpen(bool)
	pathComponent() string
	parent() FileNode
	child(name string) (FileNode, error)
	stat() (p9p.Dir, error)
}

func path(file FileNode) string {
	parent := file.parent()
	if parent == file {
		return file.pathComponent()
	} else {
		return path(file.parent()) + "/" + file.pathComponent()
	}
}

const IOUNIT uint32 = 0

func open(file FileNode, mode p9p.Flag) error {
	// TODO verify that mode is sensible . . . this should probably
	// become a FileNode method
	if file.isOpen() {
		return errors.New("file already open")
	} else {
		file.setOpen(true)
		return nil
	}
}

func clunk(file FileNode) error {
	if file.isOpen() {
		file.setOpen(false)
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
