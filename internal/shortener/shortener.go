package shortener

import (
	"bufio"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

var linksDelimiter = "|"

// Shortener makes links shorter.
type Shortener struct {
	prefix string
	links  map[string]string
	file   *os.File
}

// NewShortener creates new shortener.
func NewShortener(prefix string) Shortener {
	storagePath := os.Getenv("FILE_STORAGE_PATH")
	if storagePath == "" {
		storagePath = "links.txt"
	}

	file, err := os.OpenFile(storagePath, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}

	shortener := Shortener{
		prefix: prefix,
		links:  make(map[string]string),
		file:   file,
	}

	shortener.retrieveStoredLinks(file)

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

		s.storeLinkOnDisk(link, id)
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

func (s Shortener) storeLinkOnDisk(link string, id string) {
	newRecord := fmt.Sprintf("%s%s%s\n", link, linksDelimiter, id)

	_, err := s.file.Write([]byte(newRecord))
	if err != nil {
		panic(err)
	}
}

func (s Shortener) retrieveStoredLinks(file *os.File) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), linksDelimiter)

		if len(line) > 1 {
			s.links[line[0]] = line[1]
		}
	}
}
