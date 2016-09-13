package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"log"
)

type GistSession struct {
	client *github.Client
	fidMap map[p9p.Fid]FileNode
}

func NewGistSession(token string) *GistSession {
	var gs GistSession
	gs.client = githubClientFromToken(token)
	gs.fidMap = make(map[p9p.Fid]FileNode)
	return &gs
}

func (gs *GistSession) Auth(ctx context.Context, afid p9p.Fid, uname, aname string) (p9p.Qid, error) {
	return p9p.Qid{}, errors.New("no auth")
}

func (gs *GistSession) Attach(ctx context.Context, fid, afid p9p.Fid, uname, aname string) (p9p.Qid, error) {
	log.Println("attaching")
	rootNode := NewRootNode(gs.client)
	gs.fidMap[fid] = rootNode
	return rootNode.qid, nil
}

func (gs *GistSession) Clunk(ctx context.Context, fid p9p.Fid) error {
	if file, ok := gs.fidMap[fid]; ok {
		delete(gs.fidMap, fid)
		err := clunk(file)
		return err
	} else {
		return errors.New("don't know that fid")
	}
}

func (gs *GistSession) Remove(ctx context.Context, fid p9p.Fid) error {
	return errors.New("cant remove yet")
}

func (gs *GistSession) Write(ctx context.Context, fid p9p.Fid, p []byte, offset int64) (int, error) {
	return 0, errors.New("cant write yet")
}

func (gs *GistSession) Create(ctx context.Context, parent p9p.Fid, name string, perm uint32, mode p9p.Flag) (p9p.Qid, uint32, error) {
	return p9p.Qid{}, 0, errors.New("cant create yet")
}

func (gs *GistSession) Version() (int, string) {
	return 2048, "9p2000"
}
