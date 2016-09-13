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
		// TODO this is only correct for "directories", implement
		// something sinsible for normal files. i like having all this
		// business here though becuase it abstracts the messsy
		// encoding part
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
	} else {
		return 0, errors.New("cant find that fid")
	}
}
