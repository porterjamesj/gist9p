package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
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

func statUser(file File, user *github.User, gists []*github.Gist) p9p.Dir {
	// user updated time is most recent updated time of all the user's
	// gists. TODO using this as access time to is a bit bogo, should
	// maybe track this internally
	var times = make([]time.Time, 1)
	times[0] = user.CreatedAt.Time
	for _, gist := range gists {
		times = append(times, *gist.UpdatedAt)
	}
	modTime := maxTime(times)
	var dir = p9p.Dir{
		Type:       0,
		Dev:        0,
		Qid:        file.qid,
		Mode:       0755 | p9p.DMDIR,
		AccessTime: modTime,
		ModTime:    modTime,
		// per https://swtch.com/plan9port/man/man9/stat.html,
		// "Directories and most files representing devices have a
		// conventional length of 0. "
		Length: 0,
		Name:   file.path}
	return dir
}

func (gs *GistSession) Stat(ctx context.Context, fid p9p.Fid) (p9p.Dir, error) {
	log.Println("stating fid", fid)
	// TODO move user up to GistSession so we only get it once
	user := os.Getenv("USER")
	if file, ok := gs.store.getFid(fid); ok {
		components := strings.Split(file.path, "/")[1:]
		components = removeEmptyStrings(components)
		log.Println(fid, file, file.path)
		var dir p9p.Dir
		var err error = nil
		switch len(components) {
		case 0:
			dir = statRoot(file)
		case 1:
			uname := components[0]
			// TODO handle errors
			user, _, _ := gs.client.Users.Get(uname)
			gists, _, _ := gs.client.Gists.List(uname, nil)
			dir = statUser(file, user, gists)
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
