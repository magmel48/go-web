package shortener

import (
	"context"
	"errors"
	"fmt"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db"
	"github.com/magmel48/go-web/internal/db/links"
	"github.com/magmel48/go-web/internal/db/userlinks"
	"net/url"
)

var delimiter = "|"

// Shortener makes links shorter.
type Shortener struct {
	prefix    string
}

type UrlsMap struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewShortener creates new shortener.
func NewShortener(prefix string) Shortener {
	shortener := Shortener{
		prefix:    prefix,
	}

	db.CreateSchema()

	return shortener
}

// MakeShorter makes a link shorter.
func (s Shortener) MakeShorter(ctx context.Context, originalURL string, userID auth.UserID) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("cannot parse url")
	}

	// retrieve short link identifier
	link, err := links.FindByOriginalURL(ctx, originalURL)
	if err != nil {
		return "", err
	}

	if link == nil {
		link, err = links.Create(ctx, "", originalURL)
		if err != nil {
			return "", err
		}
	}

	// store the link for the userID if needed
	if userID != nil {
		userLink, _ := userlinks.FindByLinkID(ctx, userID, link.ID)
		if userLink == nil {
			err = userlinks.Create(ctx, userID, link.ID)
			if err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf("%s/%s", s.prefix, link.ShortID), nil
}

// RestoreLong restores short link to initial state if an info was stored before.
func (s Shortener) RestoreLong(ctx context.Context, shortID string) (string, error) {
	link, err := links.FindByShortID(ctx, shortID)
	if err != nil {
		return "", err
	}

	if link == nil {
		return "", errors.New("not found")
	}

	return link.OriginalURL, nil
}

// GetUserLinks returns all links belongs to specified userID.
func (s Shortener) GetUserLinks(ctx context.Context, userID auth.UserID) ([]UrlsMap, error) {
	if userID == nil {
		return nil, nil
	}

	userLinksList, err := userlinks.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]UrlsMap, len(userLinksList))
	for i, userLink := range userLinksList {
		result[i] = UrlsMap{
			ShortURL: fmt.Sprintf("%s/%s", s.prefix, userLink.Link.ShortID),
			OriginalURL: userLink.Link.OriginalURL,
		}
	}

	return result, nil
}
