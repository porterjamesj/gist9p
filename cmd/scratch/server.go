package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/docker/go-p9p"
	"golang.org/x/net/context"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net"
	"time"
)

func hashPath(s string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(s))
	return hash.Sum64()
}

func encodeDirs(dirs []p9p.Dir) []byte {
	codec := p9p.NewCodec()
	var buf bytes.Buffer
	for _, dir := range dirs {
		p9p.EncodeDir(codec, &buf, &dir)
	}
	data, err := ioutil.ReadAll(&buf)
	if err != nil {
		panic("failed to encode")
	}
	return data
}

func getChildren() []p9p.Dir {
	now := time.Now()
	var dirs []p9p.Dir
	for _, path := range []string{"a", "b", "c"} {
		var qid = p9p.Qid{
			Type:    p9p.QTFILE,
			Version: 0,
			Path:    hashPath(path),
		}
		var dir = p9p.Dir{
			Type: 0, Dev: 0,
			Qid:        qid,
			Mode:       0755,
			AccessTime: now,
			ModTime:    now,
			Length:     10,
			Name:       path,
			UID:        "someone",
			GID:        "someone",
			MUID:       "someone",
		}
		dirs = append(dirs, dir)
	}
	return dirs
}

func getRootInfo() (p9p.Dir, p9p.Qid) {
	now := time.Now()
	var path = "/"
	var qid = p9p.Qid{
		Type:    p9p.QTDIR,
		Version: 0,
		Path:    hashPath(path),
	}
	var dir = p9p.Dir{
		Type:       0,
		Dev:        0,
		Qid:        qid,
		Mode:       0755 | p9p.DMDIR,
		AccessTime: now,
		ModTime:    now,
		Length:     1024,
		Name:       path,
		UID:        "someone",
		GID:        "someone",
		MUID:       "someone"}
	return dir, qid
}

func main() {
	hf := p9p.HandlerFunc(func(ctx context.Context, msg p9p.Message) (p9p.Message, error) {
		var rootdir, rootQid = getRootInfo()
		fmt.Println(msg.Type())
		if msg.Type() == p9p.Tattach {
			reply := p9p.MessageRattach{rootQid}
			return reply, nil
		} else if msg.Type() == p9p.Tstat {
			return p9p.MessageRstat{rootdir}, nil
		} else if msg.Type() == p9p.Tclunk {
			return p9p.MessageRclunk{}, nil
		} else if msg.Type() == p9p.Twalk {
			fmt.Println(fmt.Sprintf("Twalk: %s", msg.(p9p.MessageTwalk)))
			var qids = []p9p.Qid{rootQid}
			fmt.Println(qids)
			return p9p.MessageRwalk{}, nil
		} else if msg.Type() == p9p.Topen {
			return p9p.MessageRopen{rootQid, 1024}, nil
		} else if msg.Type() == p9p.Tread {
			fmt.Println(fmt.Sprintf("Tread: %s", msg.(p9p.MessageTread)))
			msg := msg.(p9p.MessageTread)
			var data []byte
			if msg.Offset == 0 {
				data = encodeDirs(getChildren())
			}
			return p9p.MessageRread{data}, nil
		} else {
			return nil, errors.New(fmt.Sprintf("dont know how to handle %s", msg.Type()))
		}
	})

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal("done goofed")
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal("done goofed")
		}
		log.Println("serving connection")
		go p9p.ServeConn(context.Background(), conn, hf)
	}
}
