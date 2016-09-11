package gist9p

import (
	"errors"
	"fmt"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
)

func (gs *GistSession) Walk(ctx context.Context, fid p9p.Fid, newfid p9p.Fid, names ...string) ([]p9p.Qid, error) {
	log.Println("walking", fid, newfid, names)
	fmt.Printf("%q\n", names)
	names = removeEmptyStrings(names)
	fmt.Printf("%q\n", names)
	file, ok := gs.store.getFid(fid)
	if ok {
		var qids []p9p.Qid
		if file.path == "/" {
			if len(names) == 0 {
				//associate the file with the fid
				gs.store.fidMap[newfid] = file
				return qids, nil
			} else if len(names) == 1 {
				// name is a users, query github api
				name := names[0]
				// TODO handle errors
				user, _, _ := gs.client.Users.Get(name)
				if user != nil {
					userPath := "/" + *user.Login
					file, ok = gs.store.getPath(userPath)
					if ok {
						qids = append(qids, file.qid)
						gs.store.addFile(file, newfid)
					} else {
						userDir := NewDir(userPath)
						gs.store.addFile(userDir, newfid)
						qids = append(qids, userDir.qid)
					}
					return qids, nil
				} else {
					return nil, errors.New("user doesn't exist")
				}
			} else {
				return nil, errors.New("can't handle longer walks yet")
			}
		} else {
			// walking from non root
			return nil, errors.New("cant do non-root walks yet")
		}
	} else {
		return nil, errors.New("cant find fid to walk from")
	}
}
