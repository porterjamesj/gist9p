package gist9p

import (
	"errors"
	"fmt"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
	"regexp"
)

// yuk, i don't like this because i don't want GistSession to turn
// into a god object-y thing. i have a marvelous plan to refactor this
// but it is too small to fit in this margin
func (gs *GistSession) rootWalk(file File, names []string, newfid p9p.Fid) ([]p9p.Qid, error) {
	var qids []p9p.Qid
	switch {
	case len(names) == 0 || (len(names) == 1 && (names[0] == ".." || names[0] == ".")):
		//associate the file with the fid
		gs.store.fidMap[newfid] = file
		// qids = append(qids, file.qid)
		return qids, nil
	case len(names) == 1:
		// name is a users, query github api
		name := names[0]
		// TODO handle errors
		user, _, _ := gs.client.Users.Get(name)
		if user != nil {
			userPath := "/" + *user.Login
			userFile, ok := gs.store.getPath(userPath)
			if ok {
				qids = append(qids, userFile.qid)
				gs.store.addFile(userFile, newfid)
			} else {
				userDir := NewDir(userPath)
				gs.store.addFile(userDir, newfid)
				qids = append(qids, userDir.qid)
			}
			return qids, nil
		} else {
			return nil, errors.New("user doesn't exist")
		}
	default:
		return nil, errors.New("can't handle longer walks yet")
	}
}

func (gs *GistSession) userWalk(file File, names []string, newfid p9p.Fid) ([]p9p.Qid, error) {
	var qids []p9p.Qid
	switch {
	case len(names) == 0:
		//associate the file with the fid
		gs.store.fidMap[newfid] = file
		// qids = append(qids, file.qid)
		return qids, nil
	case len(names) == 1 && names[0] == "..":
		// walking back up to the root
		if file, ok := gs.store.getPath("/"); ok {
			gs.store.fidMap[newfid] = file
			qids = append(qids, file.qid)
			return qids, nil
		} else {
			panic("root is missing")
		}
	default:
		return nil, errors.New("don't know how to walk that yet")
	}
}

func (gs *GistSession) Walk(ctx context.Context, fid p9p.Fid, newfid p9p.Fid, names ...string) ([]p9p.Qid, error) {
	log.Println("walking", fid, newfid, names)
	fmt.Printf("%q\n", names)
	names = removeEmptyStrings(names)
	fmt.Printf("%q\n", names)
	file, ok := gs.store.getFid(fid)
	if ok {
		var qids []p9p.Qid
		var err error
		if file.path == "/" {
			qids, err = gs.rootWalk(file, names, newfid)
		} else if match, _ := regexp.MatchString(`^/[A-Za-z0-9\-]*$`, file.path); match {
			// walking from a user
			qids, err = gs.userWalk(file, names, newfid)
		} else {
			qids, err = nil, errors.New("cant do other walks yet")
		}
		return qids, err
	} else {
		return nil, errors.New("cant find fid to walk from")
	}
}
