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

func Create(ctx context.Context, shortID string, originalURL string) (*Link, bool, error) {
	if shortID == "" {
		linksCount, _ := getLinksCount(ctx)
		shortID = strconv.Itoa(linksCount)
	}

	link := Link{}
	if err := db.DB.QueryRowContext(
		ctx,
		`
			INSERT INTO "links" ("short_id", "original_url") VALUES ($1, $2)
			ON CONFLICT ("original_url") DO UPDATE SET "original_url" = "links"."original_url"
			RETURNING "id", "short_id"
		`,
		shortID,
		originalURL).Scan(&link.ID, &link.ShortID); err != nil {

		return nil, false, err
	}

	// shortID != link.ShortID if short_id`s are not the same
	return &link, shortID != link.ShortID, nil
}

func CreateBatch(ctx context.Context, originalURLs []string) ([]Link, error) {
	result := make([]Link, len(originalURLs))
	linksCount, _ := getLinksCount(ctx)

	tx, err := db.DB.Begin()
	if err != nil {
		return nil, err
	}

	defer tx.Rollback()

	insertStmt, err := tx.PrepareContext(
		ctx, `INSERT INTO "links" ("short_id", "original_url") VALUES($1, $2) RETURNING "id", "short_id"`)
	if err != nil {
		return nil, err
	}

	selectStmt, err := tx.PrepareContext(
		ctx, `SELECT "id", "short_id" FROM "links" WHERE "original_url" = $1 LIMIT 1`)
	if err != nil {
		return nil, err
	}

	txInsertStmt := tx.StmtContext(ctx, insertStmt)
	txSelectStmt := tx.StmtContext(ctx, selectStmt)

	for i, el := range originalURLs {
		link := Link{}
		rows, err := txSelectStmt.QueryContext(ctx, el)
		if err != nil {
			return nil, err
		}

		if rows.Next() {
			if err = rows.Scan(&link.ID, &link.ShortID); err != nil {
				return nil, err
			}
		} else {
			if err = txInsertStmt.QueryRowContext(ctx, strconv.Itoa(linksCount + i), el).Scan(&link.ID, &link.ShortID); err != nil {
				return nil, err
			}
		}

		result[i] = link

		if err = rows.Close(); err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return result, nil
}

func FindByShortID(ctx context.Context, shortID string) (*Link, error) {
	rows, err := db.DB.QueryContext(
		ctx, `SELECT "id", "short_id", "original_url" FROM "links" WHERE "short_id" = $1 LIMIT 1`, shortID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	if rows.Next() {
		link := Link{}
		if err := rows.Scan(&link.ID, &link.ShortID, &link.OriginalURL); err != nil {
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

	if rows.Next() {
		link := Link{}
		if err := rows.Scan(&link.ID, &link.ShortID, &link.OriginalURL); err != nil {
			return nil, err
		}

		return &link, nil
	}

	return nil, nil
}

func getLinksCount(ctx context.Context) (int, error) {
	count := 0

	rows, err := db.DB.QueryContext(ctx, `SELECT COUNT(*) FROM "links"`)
	if err != nil {
		return 0, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}

		count = count + 1
	}

	return count, nil
}
