package userlinks

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db/links"
	"strings"
)

// PostgresRepository is implementation of abstract Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository returns new PostgresRepository for working with user links.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create creates new record (relation) by user and link identifier.
func (repository *PostgresRepository) Create(ctx context.Context, userID auth.UserID, linkID int) error {
	_, err := repository.db.ExecContext(
		ctx, `INSERT INTO "user_links" ("user_id", "link_id") VALUES ($1, $2)`, *userID, linkID)

	return err
}

// List returns list of user links.
func (repository *PostgresRepository) List(ctx context.Context, userID auth.UserID) ([]UserLink, error) {
	rows, err := repository.db.QueryContext(
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

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return result, nil
}

// FindByLinkID finds a link by user and link identifier.
func (repository *PostgresRepository) FindByLinkID(ctx context.Context, userID auth.UserID, linkID int) (*UserLink, error) {
	rows, err := repository.db.QueryContext(
		ctx,
		`SELECT "id", "user_id", "link_id" FROM "user_links" WHERE "user_id" = $1 AND "link_id" = $2 LIMIT 1`,
		*userID,
		linkID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	userLink := UserLink{}
	if rows.Next() {
		err := rows.Scan(&userLink.ID, &userLink.UserID, &userLink.LinkID)
		if err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if userLink.ID != 0 {
		return &userLink, nil
	}

	return nil, nil
}

// DeleteLinks deletes user links by batches with many links inside.
func (repository *PostgresRepository) DeleteLinks(ctx context.Context, deleteQueryItems []DeleteQueryItem) error {
	query := `
		UPDATE "links" AS l
		SET "is_deleted" = true
		FROM "user_links" AS ul WHERE
	`

	clauses := make([]string, len(deleteQueryItems))
	for i := range deleteQueryItems {
		clause := fmt.Sprintf(`l."short_id" = ANY ($%d) AND l."id" = ul."link_id" AND ul."user_id" = $%d`, 2*i+1, 2*i+2)
		clauses[i] = clause
	}

	query += strings.Join(clauses, " OR ")

	args := make([]interface{}, len(deleteQueryItems)*2)
	for i, el := range deleteQueryItems {
		args[2*i] = el.ShortIDs
		args[2*i+1] = *el.UserID
	}

	_, err := repository.db.ExecContext(ctx, query, args...)

	return err
}
