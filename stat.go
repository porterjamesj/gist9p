package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
	"os"
	"strings"
)

func statFile(node Node) (p9p.Dir, error) {
	components := strings.Split(path(node), "/")[1:]
	components = removeEmptyStrings(components)
	dir, err := node.Stat()
	dir.Qid = node.Qid()
	dir.Name = node.PathComponent()
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
	// TODO this is totally bogo, we at the very least should:
	// 1. be willing to change the name
	// 2. intelligently log everything else the client was trying to change
	//
	// There's a semantic mismatch between the semantics of the FUSE
	// API and 9p here. In FUSE, it seems to be the case that you can
	// Setattr a file's length to 0 as a means of truncating it
	// (https://www.cs.columbia.edu/~du/ds/homework/2013/09/26/Homework3/#setattr-open-write-and-read),
	// whereas 9p specifies that the length of a file cannot be set
	// (https://swtch.com/plan9port/man/man9/stat.html).
	//
	// I guess the way to truncate things in plan9 is by opening with
	// O_TRUNC? I kind of like the FUSE API better here since you can
	// truncate more flexibly. In any case it seems like we have to
	// implement truncation a la FUSE in order to work reasonably well
	// with 9pfuse. This might be harder than anticipated since:
	//
	// 1. We need a way of detecting if the client is only trying to
	// set one attribute, and 2. there appears to be some sort of
	// encoding issue, as far as I can tell 9pFUSE encodse the Dir
	// entry such that it doesn't make it through to here
	// correctly. Sigh.
	//
	// I am very curious to see how the linux kernel implentation
	// handles this situation.
	log.Println("supposed new stat:", dir)
	return nil
}
