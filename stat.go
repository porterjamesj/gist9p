package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
)

func statFile(file FileNode) (p9p.Dir, error) {
	components := strings.Split(path(file), "/")[1:]
	components = removeEmptyStrings(components)
	dir, err := file.Stat()
	dir.Qid = file.Qid()
	dir.Name = file.PathComponent()
	dir.Type = 0
	dir.Dev = 0
	// TODO move user up to GistSession so we only get it once
	user := os.Getenv("USER")
	dir.MUID = user
	dir.UID = user
	dir.GID = user
	return dir, err
}

func (gs *GistSession) Stat(ctx context.Context, fid p9p.Fid) (p9p.Dir, error) {
	log.Println("stating fid", fid)
	if file, ok := gs.fidMap[fid]; ok {
		dir, err := statFile(file)
		return dir, err
	} else {
		return p9p.Dir{}, errors.New("fid not found")
	}
}

func (gs *GistSession) WStat(ctx context.Context, fid p9p.Fid, dir p9p.Dir) error {
	return errors.New("cant write stat info yet")
}
