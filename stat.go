package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
)

func (gs *GistSession) Stat(ctx context.Context, fid p9p.Fid) (p9p.Dir, error) {
	log.Println("stating fid", fid)
	if file, ok := gs.fidMap[fid]; ok {
		components := strings.Split(path(file), "/")[1:]
		components = removeEmptyStrings(components)
		log.Println(fid, file, path(file))
		dir, err := file.stat()
		dir.Qid = file.getQid()
		dir.Name = path(file)
		dir.Type = 0
		dir.Dev = 0
		// TODO move user up to GistSession so we only get it once
		user := os.Getenv("USER")
		dir.MUID = user
		dir.UID = user
		dir.GID = user
		return dir, err
	} else {
		return p9p.Dir{}, errors.New("fid not found")
	}
}

func (gs *GistSession) WStat(ctx context.Context, fid p9p.Fid, dir p9p.Dir) error {
	return errors.New("cant write stat info yet")
}
