package gist9p

import (
	"bytes"
	"github.com/docker/go-p9p"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"hash/fnv"
	"io/ioutil"
	"time"
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

func removeEmptyStrings(strings []string) []string {
	var cleaned []string
	for _, s := range strings {
		if s != "" {
			cleaned = append(cleaned, s)
		}
	}
	return cleaned
}

func maxTime(times []time.Time) time.Time {
	// TODO i don't like returning this totally bogus time
	max := time.Time{}
	for _, t := range times {
		if t.After(max) {
			max = t
		}
	}
	return max
}

func encodeDirs(dirs []p9p.Dir) ([]byte, error) {
	codec := p9p.NewCodec()
	var buf bytes.Buffer
	for _, dir := range dirs {
		p9p.EncodeDir(codec, &buf, &dir)
	}
	data, err := ioutil.ReadAll(&buf)
	if err != nil {
		return nil, err
	}
	return data, nil
}
