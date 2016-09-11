package gist9p

import (
	"errors"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"hash/fnv"
	"log"
)

func githubClientFromToken(token string) *github.Client {
	// https://godoc.org/github.com/google/go-github/github#hdr-Authentication
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})
	tc := oauth2.NewClient(oauth2.NoContext, ts)
	return github.NewClient(tc)
}

func hashPath(s string) uint64 {
	hash := fnv.New64()
	hash.Write([]byte(s))
	return hash.Sum64()
}

// TODO move this into file store
func getRootFile() File {
	var path = "/"
	var qid = p9p.Qid{
		Type:    p9p.QTDIR,
		Version: 0,
		Path:    hashPath(path),
	}
	return File{
		qid:  qid,
		path: path,
	}
}

type GistSession struct {
	client *github.Client
	// TODO pull this out into a "middleware" session, that wraps
	// an inner session and does fid / qid mapping bookkeeping by
	// keeping a map from fid -> qid, and passing the relevant qid
	// into requests via the context k/v store. is this a
	// reasonable / idiomatic usage of context?
	store *FileStore
}

func NewGistSession(token string) *GistSession {
	var gs GistSession
	gs.client = githubClientFromToken(token)
	gs.store = NewFileStore()
	return &gs
}

func (gs *GistSession) Auth(ctx context.Context, afid p9p.Fid, uname, aname string) (p9p.Qid, error) {
	return p9p.Qid{}, errors.New("no auth")
}

func (gs *GistSession) Attach(ctx context.Context, fid, afid p9p.Fid, uname, aname string) (p9p.Qid, error) {
	log.Println("attaching")
	rootFile := getRootFile()
	gs.store.addFile(rootFile, fid)
	return rootFile.qid, nil
}

func (gs *GistSession) Clunk(ctx context.Context, fid p9p.Fid) error {
	return nil
}

func (gs *GistSession) Remove(ctx context.Context, fid p9p.Fid) error {
	return errors.New("cant remove yet")
}

func (gs *GistSession) Read(ctx context.Context, fid p9p.Fid, p []byte, offset int64) (int, error) {
	return 0, errors.New("cant read yet")
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
