package gist9p

import (
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"hash/fnv"
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
