package links

import (
	"context"
	"errors"
)

// Link is representing database table and a link DTO at the same time.
type Link struct {
	ID          int
	ShortID     string
	OriginalURL string
	IsDeleted   bool
}

// Repository is common interface for a work with links implementation.
//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, shortID string, originalURL string) (*Link, error)
	CreateBatch(ctx context.Context, originalURLs []string) ([]Link, error)
	FindByShortID(ctx context.Context, shortID string) (*Link, error)
}

// ErrConflict is using for notifying clients about a conflict with shorter link identifiers. Usually it means
// the link was made already shorter, but by another user.
var ErrConflict = errors.New("conflict")
