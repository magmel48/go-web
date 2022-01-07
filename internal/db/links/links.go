package links

import (
	"context"
	"github.com/magmel48/go-web/internal/db"
	"strconv"
)

type Link struct {
	ID          int
	ShortID     string
	OriginalURL string
}

func Create(ctx context.Context, shortID string, originalURL string) (*Link, error) {
	if shortID == "" {
		count := 0

		rows, err := db.DB.QueryContext(ctx, `SELECT COUNT(*) FROM "links"`)
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		for rows.Next() {
			rows.Scan(&count)
			count = count + 1
		}

		shortID = strconv.Itoa(count)
	}

	id := 0
	err := db.DB.QueryRowContext(
		ctx, `INSERT INTO "links" ("short_id", "original_url") VALUES ($1, $2) RETURNING id`, shortID, originalURL).Scan(&id)
	if err != nil {
		return nil, err
	}

	return &Link{ID: id, ShortID: shortID, OriginalURL: originalURL}, nil
}

func FindByShortID(ctx context.Context, shortID string) (*Link, error) {
	rows, err := db.DB.QueryContext(
		ctx, `SELECT "id", "short_id", "original_url" FROM "links" WHERE "short_id" = $1 LIMIT 1`, shortID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		link := Link{}
		err := rows.Scan(&link.ID, &link.ShortID, &link.OriginalURL)
		if err != nil {
			return nil, err
		}

		return &link, nil
	}

	return nil, nil
}

func FindByOriginalURL(ctx context.Context, originalURL string) (*Link, error) {
	rows, err := db.DB.QueryContext(
		ctx, `SELECT "id", "short_id", "original_url" FROM "links" WHERE "original_url" = $1 LIMIT 1`, originalURL)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		link := Link{}
		err := rows.Scan(&link.ID, &link.ShortID, &link.OriginalURL)
		if err != nil {
			return nil, err
		}

		return &link, nil
	}

	return nil, nil
}
