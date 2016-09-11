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
	file, ok := gs.fidMap[fid]
	if ok {
		var qids []p9p.Qid
		var err error
		curr := file
		for _, name := range names {
			if name == "." {
				// curr remains the same
			} else if name == ".." {
				curr = curr.parent()
			} else {
				curr, err = curr.child(name)
				if err != nil {
					// TODO do we return the qids that *have* been
					// successfully walked? or do we return an empty
					// list? should read over docs
					break
				}
			}
			qids = append(qids, curr.getQid())
		}
		gs.fidMap[newfid] = curr
		return qids, err
	} else {
		return nil, errors.New("cant find fid to walk from")
	}
}
