package tin

import (
	"log"
	"regexp"
	"strings"
)

type path struct {
	url     string
	params  map[string]int
	pattern *regexp.Regexp
}

func extractPath(url string) *path {

	parts := make([]string, 0)
	params := make(map[string]int)

	for i, part := range strings.Split(url, "/") {

		if len(part) > 0 && part[0] == ':' {
			params[part[1:]] = i
			parts = append(parts, "[^/]*")

		} else {
			parts = append(parts, part)
		}
	}

	url = strings.Join(parts, "/")
	pattern, err := regexp.Compile("^" + url + "$")
	if err != nil {
		log.Fatal("Invalid path", url, err)
	}

	if len(params) == 0 {
		params = nil
	}

	return &path{url, params, pattern}

}

func (path *path) match(url string) bool {
	return path.pattern.MatchString(url)
}

func (path *path) param(url string, key string) string {
	if idx, has := path.params[key]; has {
		url := strings.Split(url, "/")
		return url[idx]
	}
	panic("Invalid parameter")
}
