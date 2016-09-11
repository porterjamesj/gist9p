package gist9p

import (
	"time"
)

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
