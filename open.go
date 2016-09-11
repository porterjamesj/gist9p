package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
)

// TODO figure out what the actual correct value is here
const IOUNIT uint32 = 0

func (gs *GistSession) Open(ctx context.Context, fid p9p.Fid, mode p9p.Flag) (p9p.Qid, uint32, error) {
	log.Println("opening fid", fid)
	if file, ok := gs.store.getFid(fid); ok {
		if mode == p9p.OREAD {
			// TODO error if file is already open?
			file.open = true
			return file.qid, IOUNIT, nil
		} else {
			return p9p.Qid{}, IOUNIT, errors.New("cant open for non-reading yet")
		}
	} else {
		return p9p.Qid{}, IOUNIT, errors.New("cant find that fid")
	}
}
