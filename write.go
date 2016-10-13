package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
)

func (gs *GistSession) Write(ctx context.Context, fid p9p.Fid, p []byte, offset int64) (int, error) {
	log.Println("write fid", fid, "offset", offset, "count", len(p))
	log.Println("write data", string(p))
	if file, ok := gs.fidMap[fid]; ok {
		written, err := file.Write(p, offset)
		if written != len(p) {
			return 0, errors.New("failed to write entire request")
		}
		return written, err
	} else {
		return 0, errors.New("can't find file")
	}
}
