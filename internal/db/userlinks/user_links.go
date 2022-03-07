package userlinks

import (
	"context"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db/links"
)

// UserLink is representing database table and a user link DTO at the same time.
type UserLink struct {
	ID     int
	UserID auth.UserID
	LinkID int
	Link   links.Link
}

// DeleteQueryItem is specifying deleting intention from the UserID.
type DeleteQueryItem struct {
	UserID   auth.UserID
	ShortIDs []string
}

// Repository is common interface for a work with user links implementation.
//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, userID auth.UserID, linkID int) error
	List(ctx context.Context, userID auth.UserID) ([]UserLink, error)
	FindByLinkID(ctx context.Context, userID auth.UserID, linkID int) (*UserLink, error)
	DeleteLinks(ctx context.Context, deleteQueryItems []DeleteQueryItem) error
}
