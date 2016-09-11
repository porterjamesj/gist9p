package gist9p

func removeEmptyStrings(strings []string) []string {
	var cleaned []string
	for _, s := range strings {
		if s != "" {
			cleaned = append(cleaned, s)
		}
	}
	return cleaned
}
