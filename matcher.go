package apicore

import "strings"

type Matcher interface {
	Match(url string, method string) bool
}

type matcher string

func (m matcher) Match(url string, method string) bool {
	path, methods := m.parse()
	url = strings.Split(url, "?")[0]
	if url == path {
		for _, m := range methods {
			if m == "ALL" {
				return true
			}
			if m != method {
				continue
			}
			return true
		}
	}
	return false
}

func (m matcher) parse() (url string, path []string) {
	args := strings.Split(string(m), "|")
	return args[0], args[1:]
}

func NewMatcher(path string, methods ...string) Matcher {
	if len(methods) == 0 {
		methods = []string{"ALL"}
	}

	return matcher(strings.Join(append([]string{path}, methods...), "|"))
}
