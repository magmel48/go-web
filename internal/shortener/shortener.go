package shortener

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magmel48/go-web/internal/auth"
	"log"
	"net/url"
)

var delimiter = "|"

// Shortener makes links shorter.
type Shortener struct {
	prefix    string
	links     map[string]string
	userLinks map[string][]string
	backup    Backup
}

type UrlsMap struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewShortener creates new shortener.
func NewShortener(prefix string, store Backup) Shortener {
	shortener := Shortener{
		prefix:    prefix,
		links:     make(map[string]string),
		userLinks: make(map[string][]string),
		backup:    store,
	}

	shortener.retrieveStoredData()

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
		if _, ok := s.userLinks[userID.String()]; !ok {
			s.userLinks[userID.String()] = make([]string, 0)
		}

		found := false
		for _, el := range s.userLinks[userID.String()] {
			if el == linkID {
				found = true
				break
			}
		}

		if !found {
			s.userLinks[userID.String()] = append(s.userLinks[userID.String()], linkID)
			s.storeUserLinkData(userID)
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
	if userID == nil {
		return nil
	}

	if userLinks, ok := s.userLinks[userID.String()]; !ok {
		return nil
	} else {
		result := make([]UrlsMap, len(userLinks))

		for i, linkID := range userLinks {
			longLink, _ := s.RestoreLong(linkID)
			result[i] = UrlsMap{ShortURL: fmt.Sprintf("%s/%s", s.prefix, linkID), OriginalURL: longLink}
		}

		return result
	}
}

func (s Shortener) storeLink(link string, linkID string) {
	record := fmt.Sprintf("%s%s%s\n", link, delimiter, linkID)
	err := s.backup.Append(record)

	if err != nil {
		log.Printf("not able to store link %s", link)
	}
}

func (s Shortener) retrieveStoredData() {
	data := s.backup.ReadAll()

	for k, v := range data {
		_, err := url.ParseRequestURI(k)

		if err != nil {
			IDs := make([]string, 0)
			err := json.Unmarshal([]byte(v), &IDs)
			if err != nil {
				log.Println("user links unmarshal error", err)
			}

			s.userLinks[k] = IDs
		} else {
			s.links[k] = v
		}
	}

	log.Println("restored", s.links, s.userLinks)
}

func (s Shortener) storeUserLinkData(userID auth.UserID) {
	marshaledIDs, _ := json.Marshal(s.userLinks[userID.String()])

	record := fmt.Sprintf(
		"%s%s%s\n", userID.String(), delimiter,  string(marshaledIDs))
	err := s.backup.Append(record)

	if err != nil {
		log.Printf("not able to store data for user %s", userID.String())
	}
}
