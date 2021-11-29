package shortener

import (
	"errors"
	"fmt"
	"net/url"
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

func (s Shortener) MakeShorter(link string) (string, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return "", errors.New("cannot parse url")
	}

	id := ""

	if storedLink, ok := s.links[link]; ok {
		id = storedLink
	} else {
		id = fmt.Sprintf("%d", len(s.links)+1)
		s.links[link] = id
	}

	return fmt.Sprintf("%s/%s", s.prefix, id), nil
}

func (s Shortener) RestoreLong(id string) (string, error) {
	for k, v := range s.links {
		if v == id {
			return k, nil
		}
	}

	return "", errors.New("not found")
}
