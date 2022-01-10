package links

import (
	"context"
	"errors"
)

type Link struct {
	ID          int
	ShortID     string
	OriginalURL string
}

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, shortID string, originalURL string) (*Link, error)
	CreateBatch(ctx context.Context, originalURLs []string) ([]Link, error)
	FindByShortID(ctx context.Context, shortID string) (*Link, error)
}

var ErrConflict = errors.New("conflict")
