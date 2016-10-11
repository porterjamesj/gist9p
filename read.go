package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"log"
)

func (gs *GistSession) Read(ctx context.Context, fid p9p.Fid, p []byte, offset int64) (int, error) {
	log.Println("read fid offset", fid, offset)
	if file, ok := gs.fidMap[fid]; ok {
		var qid = file.Qid()
		if qid.Type == p9p.QTDIR {
			// TODO figure out how to use the Readdir type from p9p,
			// which has a much more robust implementation of this
			// that's a bit hard to figure out how to fit in here
			// because its stateful, so we need to somehow attach it
			// to the interface
			children, err := file.Children()
			var dirs []p9p.Dir
			for _, child := range children {
				dir, err := statFile(child)
				if err != nil {
					return 0, err
				}
				dirs = append(dirs, dir)
			}
			bytes, err := encodeDirs(dirs)
			if err != nil {
				return 0, err
			}
			if int(offset) > len(bytes) {
				return 0, nil
			} else {
				copy(p, bytes[offset:])
				return len(bytes) - int(offset), nil
			}
		} else if qid.Type == p9p.QTFILE {
			read, err := file.Read(p, offset)
			return read, err
		} else {
			return 0, errors.New("don't how to read that QType")
		}
	} else {
		return 0, errors.New("cant find that fid")
	}
}
