package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
	"time"
)

func statRoot(file File) p9p.Dir {
	// TODO think of something more principled to do about access
	// time?
	now := time.Now()
	var dir = p9p.Dir{
		Type:       0,
		Dev:        0,
		Qid:        file.qid,
		Mode:       0755 | p9p.DMDIR,
		AccessTime: now,
		ModTime:    now,
		Length:     0,
		Name:       file.path}
	return dir
}

func (gs *GistSession) Stat(ctx context.Context, fid p9p.Fid) (p9p.Dir, error) {
	log.Println("stating")
	// TODO move user up to gist session
	user := os.Getenv("USER")
	if file, ok := gs.store.getFid(fid); ok {
		components := strings.Split(file.path, "/")[1:]
		components = removeEmptyStrings(components)
		log.Println(fid, file, file.path)
		var dir p9p.Dir
		var err error
		err = nil
		switch len(components) {
		case 0:
			dir = statRoot(file)
		default:
			dir = p9p.Dir{}
			err = errors.New("cant stat that yet")
		}
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
