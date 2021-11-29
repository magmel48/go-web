package shortener

import (
	"errors"
	"fmt"
)

// Shortener makes links shorter.
type Shortener struct {
	prefix string
	links  map[string]string
}

// NewShortener creates new shortener.
func NewShortener(prefix string) Shortener {
	return Shortener{
		prefix: prefix,
		links:  make(map[string]string),
	}
}

func (s Shortener) MakeShorter(url string) string {
	id := ""

	if link, ok := s.links[url]; ok {
		id = link
	} else {
		id = fmt.Sprintf("%d", len(s.links)+1)
		s.links[url] = id
	}

	return fmt.Sprintf("%s/%s", s.prefix, id)
}

func (s Shortener) RestoreLong(id string) (string, error) {
	for k, v := range s.links {
		if v == id {
			return k, nil
		}
	}

	return "", errors.New("not found")
}
