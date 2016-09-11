package gist9p

import (
	"github.com/docker/go-p9p"
)

type File struct {
	qid  p9p.Qid
	path string
	open bool
}

func NewFileOfType(path string, qtype p9p.QType) File {
	return File{
		qid: p9p.Qid{
			Type:    qtype,
			Version: 0,
			Path:    hashPath(path),
		},
		path: path,
		open: false,
	}
}

func NewDir(path string) File {
	return NewFileOfType(path, p9p.QTDIR)
}

func NewFile(path string) File {
	return NewFileOfType(path, p9p.QTFILE)
}

// stores Files, lets us access them by fid or string path
type FileStore struct {
	fidMap  map[p9p.Fid]File
	pathMap map[string]File
}

func (store *FileStore) getFid(fid p9p.Fid) (File, bool) {
	file, ok := store.fidMap[fid]
	return file, ok
}

func (store *FileStore) getPath(path string) (File, bool) {
	file, ok := store.pathMap[path]
	return file, ok
}

func (store *FileStore) addFile(file File, fid p9p.Fid) {
	store.pathMap[file.path] = file
	store.fidMap[fid] = file
}

func NewFileStore() *FileStore {
	var fs FileStore
	fs.fidMap = make(map[p9p.Fid]File)
	fs.pathMap = make(map[string]File)
	return &fs
}
