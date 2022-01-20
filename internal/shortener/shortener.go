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

var ErrDeleted = errors.New("the link is deleted")

// Shortener makes links shorter.
type Shortener struct {
	ctx                 context.Context
	prefix              string
	database            db.DB
	linksRepository     links.Repository
	userLinksRepository userlinks.Repository
}

type UrlsMap struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// NewShortener creates new shortener.
func NewShortener(ctx context.Context, prefix string, database db.DB) Shortener {
	shortener := Shortener{
		ctx:                 ctx,
		prefix:              prefix,
		database:            database,
		linksRepository:     links.NewPostgresRepository(database.Instance()),
		userLinksRepository: userlinks.NewPostgresRepository(database.Instance()),
	}

	return shortener
}

func (s Shortener) IsStorageAvailable(ctx context.Context) bool {
	return s.database.CheckConnection(ctx)
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
func (s Shortener) MakeShorter(ctx context.Context, originalURL string, userID auth.UserID) (string, error) {
	_, err := url.ParseRequestURI(originalURL)
	if err != nil {
		return "", errors.New("cannot parse url")
	}

	link, err := s.linksRepository.Create(ctx, "", originalURL)
	if err != nil {
		if !errors.Is(err, links.ErrConflict) {
			return "", err
		}
	}

	// store the link for the userID if needed
	if userID != nil {
		userLink, _ := s.userLinksRepository.FindByLinkID(ctx, userID, link.ID)
		if userLink == nil {
			if err := s.userLinksRepository.Create(ctx, userID, link.ID); err != nil {
				return "", err
			}
		}
	}

	return fmt.Sprintf("%s/%s", s.prefix, link.ShortID), err
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

	if link.IsDeleted {
		return link.OriginalURL, ErrDeleted
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

func (s Shortener) DeleteURLs(userID auth.UserID, IDs []string) error {
	// TODO

	return nil
}
