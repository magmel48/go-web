package shortener

import (
	"errors"
	"fmt"
	"github.com/magmel48/go-web/internal/auth"
	"log"
	"net/url"
)

var linksDelimiter = "|"

// Shortener makes links shorter.
type Shortener struct {
	prefix    string
	links     map[string]string
	userLinks map[auth.UserID][]string
	backup    Backup
}

type UrlsMap struct {
	ShortUrl    string `json:"short_url"`
	OriginalUrl string `json:"original_url"`
}

// NewShortener creates new shortener.
func NewShortener(prefix string, store Backup) Shortener {
	shortener := Shortener{
		prefix:    prefix,
		links:     make(map[string]string),
		userLinks: make(map[auth.UserID][]string),
		backup:    store,
	}

	shortener.retrieveStoredLinks()

	return shortener
}

// MakeShorter makes a link shorter.
func (s Shortener) MakeShorter(link string, userID auth.UserID) (string, error) {
	_, err := url.ParseRequestURI(link)
	if err != nil {
		return "", errors.New("cannot parse url")
	}

	// retrieve short link identifier
	linkID := ""
	if storedLink, ok := s.links[link]; ok {
		linkID = storedLink
	} else {
		linkID = fmt.Sprintf("%d", len(s.links)+1)

		s.links[link] = linkID
		s.storeLink(link, linkID)
	}

	// store the link for the userID if needed
	if userID != nil {
		if _, ok := s.userLinks[userID]; !ok {
			s.userLinks[userID] = make([]string, 0)
		}

		found := false
		for _, el := range s.userLinks[userID] {
			if el == linkID {
				found = true
				break
			}
		}

		if !found {
			s.userLinks[userID] = append(s.userLinks[userID], linkID)
		}
	}

	return fmt.Sprintf("%s/%s", s.prefix, linkID), nil
}

// RestoreLong restores short link to initial state if an info was stored before.
func (s Shortener) RestoreLong(linkID string) (string, error) {
	for k, v := range s.links {
		if v == linkID {
			return k, nil
		}
	}

	return "", errors.New("not found")
}

// GetUserLinks returns all links belongs to specified userID.
func (s Shortener) GetUserLinks(userID auth.UserID) []UrlsMap {
	if userLinks, ok := s.userLinks[userID]; !ok {
		return nil
	} else {
		result := make([]UrlsMap, len(userLinks))

		for i, shortLink := range userLinks {
			longLink, _ := s.RestoreLong(shortLink)
			result[i] = UrlsMap{ShortUrl: shortLink, OriginalUrl: longLink}
		}

		return result
	}
}

func (s Shortener) storeLink(link string, linkID string) {
	record := fmt.Sprintf("%s%s%s\n", link, linksDelimiter, linkID)
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
