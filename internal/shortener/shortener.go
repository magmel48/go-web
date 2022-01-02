package shortener

import (
	"errors"
	"fmt"
	"log"
	"net/url"
)

var linksDelimiter = "|"

// Shortener makes links shorter.
type Shortener struct {
	prefix string
	links  map[string]string
	backup Backup
}

type UrlsMap struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

// NewShortener creates new shortener.
func NewShortener(prefix string, store Backup) Shortener {
	shortener := Shortener{
		prefix: prefix,
		links:  make(map[string]string),
		backup: store,
	}

	shortener.retrieveStoredLinks()

	return shortener
}

// MakeShorter makes a link shorter.
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

		s.storeLink(link, id)
	}

	return fmt.Sprintf("%s/%s", s.prefix, id), nil
}

// RestoreLong restores short link to initial state if an info was stored before.
func (s Shortener) RestoreLong(id string) (string, error) {
	for k, v := range s.links {
		if v == id {
			return k, nil
		}
	}

	return "", errors.New("not found")
}

func (s Shortener) storeLink(link string, id string) {
	record := fmt.Sprintf("%s%s%s\n", link, linksDelimiter, id)
	err := s.backup.Append(record)

	if err != nil {
		log.Printf("not able to store link %s", link)
	}
}

func (s Shortener) retrieveStoredLinks() {
	links := s.backup.ReadAll()
	for k, v := range links {
		s.links[k] = v
	}
}
