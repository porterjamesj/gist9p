package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
)

func (gs *GistSession) Open(ctx context.Context, fid p9p.Fid, mode p9p.Flag) (p9p.Qid, uint32, error) {
	log.Println("opening fid", fid)
	if file, ok := gs.fidMap[fid]; ok {
		err := open(file, mode)
		return file.getQid(), IOUNIT, err
	} else {
		return p9p.Qid{}, IOUNIT, errors.New("cant find that fid")
	}
}
