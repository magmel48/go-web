package userlinks

import (
	"context"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db"
	"github.com/magmel48/go-web/internal/db/links"
)

type UserLink struct {
	ID     int
	UserID auth.UserID
	LinkID int
	Link   links.Link
}

func Create(ctx context.Context, userID auth.UserID, linkID int) error {
	_, err := db.DB.ExecContext(
		ctx, `INSERT INTO "user_links" ("user_id", "link_id") VALUES ($1, $2)`, *userID, linkID)

	return err
}

func List(ctx context.Context, userID auth.UserID) ([]UserLink, error) {
	rows, err := db.DB.QueryContext(
		ctx,
		`SELECT l."short_id", l."original_url" FROM "user_links" AS ul JOIN "links" as l ON ul."link_id" = l."id" WHERE ul."user_id" = $1`,
		*userID)
	if err != nil {
		return nil, err
	}

	result := make([]UserLink, 0)

	for rows.Next() {
		link := links.Link{}
		err := rows.Scan(&link.ShortID, &link.OriginalURL)
		if err != nil {
			return nil, err
		}

		result = append(result, UserLink{UserID: userID, Link: link})
	}

	return result, nil
}

func FindByLinkID(ctx context.Context, userID auth.UserID, linkID int) (*UserLink, error) {
	rows, err := db.DB.QueryContext(
		ctx,
		`SELECT "id", "user_id", "link_id" FROM "user_links" WHERE "user_id" = $1 AND "link_id" = $2`,
		*userID,
		linkID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		userLink := UserLink{}
		err := rows.Scan(&userLink.ID, &userLink.UserID, &userLink.LinkID)
		if err != nil {
			return nil, err
		}

		return &userLink, nil
	}

	return nil, nil
}
