package links

import (
	"context"
)

type Link struct {
	ID          int
	ShortID     string
	OriginalURL string
}

type Repository interface {
	Create(ctx context.Context, shortID string, originalURL string) (*Link, bool, error)
	CreateBatch(ctx context.Context, originalURLs []string) ([]Link, error)
	FindByShortID(ctx context.Context, shortID string) (*Link, error)
}
