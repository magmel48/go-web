package userlinks

import (
	"context"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db/links"
)

type UserLink struct {
	ID     int
	UserID auth.UserID
	LinkID int
	Link   links.Link
}

//go:generate mockery --name=Repository
type Repository interface {
	Create(ctx context.Context, userID auth.UserID, linkID int) error
	List(ctx context.Context, userID auth.UserID) ([]UserLink, error)
	FindByLinkID(ctx context.Context, userID auth.UserID, linkID int) (*UserLink, error)
}
