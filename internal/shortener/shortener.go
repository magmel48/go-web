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

// Shortener makes links shorter.
type Shortener struct {
	prefix              string
	linksRepository     links.Repository
	userLinksRepository userlinks.Repository
}

type UrlsMap struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewShortener creates new shortener.
func NewShortener(prefix string) Shortener {
	if err := db.CreateSchema(); err != nil {
		panic(err)
	}

	shortener := Shortener{
		prefix: prefix,
		linksRepository: links.NewPostgresRepository(db.DB),
		userLinksRepository: userlinks.NewPostgresRepository(db.DB),
	}

	return shortener
}

func (s Shortener) MakeShorterBatch(ctx context.Context, originalURLs []string) ([]string, error) {
	linkRecords, err := s.linksRepository.CreateBatch(ctx, originalURLs)
	if err != nil {
		return nil, err
	}

	result := make([]string, len(originalURLs))
	for i, link := range linkRecords {
		result[i] = fmt.Sprintf("%s/%s", s.prefix, link.ShortID)
	}

	return result, nil
}

// MakeShorter makes a link shorter.
func (s Shortener) MakeShorter(ctx context.Context, originalURL string, userID auth.UserID) (string, bool, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", false, errors.New("cannot parse url")
	}

	link, isDuplicated, err := s.linksRepository.Create(ctx, "", originalURL)
	if err != nil {
		return "", false, err
	}

	// store the link for the userID if needed
	if userID != nil {
		userLink, _ := s.userLinksRepository.FindByLinkID(ctx, userID, link.ID)
		if userLink == nil {
			if err = s.userLinksRepository.Create(ctx, userID, link.ID); err != nil {
				return "", false, err
			}
		}
	}

	return fmt.Sprintf("%s/%s", s.prefix, link.ShortID), isDuplicated, nil
}

// RestoreLong restores short link to initial state if an info was stored before.
func (s Shortener) RestoreLong(ctx context.Context, shortID string) (string, error) {
	link, err := s.linksRepository.FindByShortID(ctx, shortID)
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

	userLinksList, err := s.userLinksRepository.List(ctx, userID)
	if err != nil {
		return nil, err
	}

	result := make([]UrlsMap, len(userLinksList))
	for i, userLink := range userLinksList {
		result[i] = UrlsMap{
			ShortURL:    fmt.Sprintf("%s/%s", s.prefix, userLink.Link.ShortID),
			OriginalURL: userLink.Link.OriginalURL,
		}
	}

	return result, nil
}
