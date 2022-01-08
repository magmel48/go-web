package links

import (
	"context"
	"database/sql"
	"strconv"
)

// PostgresRepository is implementation of abstract Repository.
type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (repository *PostgresRepository) Create(ctx context.Context, shortID string, originalURL string) (*Link, bool, error) {
	if shortID == "" {
		linksCount, _ := repository.getNextShortID(ctx)
		shortID = strconv.Itoa(linksCount)
	}

	link := Link{}
	if err := repository.db.QueryRowContext(
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

func (repository *PostgresRepository) CreateBatch(ctx context.Context, originalURLs []string) ([]Link, error) {
	result := make([]Link, len(originalURLs))
	linksCount, _ := repository.getNextShortID(ctx)

	tx, err := repository.db.Begin()
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
			if err = txInsertStmt.QueryRowContext(ctx, strconv.Itoa(linksCount+i), el).Scan(&link.ID, &link.ShortID); err != nil {
				return nil, err
			}
		}

		err = rows.Err()
		if err != nil {
			return nil, err
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

func (repository *PostgresRepository) FindByShortID(ctx context.Context, shortID string) (*Link, error) {
	rows, err := repository.db.QueryContext(
		ctx, `SELECT "id", "short_id", "original_url" FROM "links" WHERE "short_id" = $1 LIMIT 1`, shortID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	link := Link{}
	if rows.Next() {
		if err := rows.Scan(&link.ID, &link.ShortID, &link.OriginalURL); err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return &link, nil
}

func (repository *PostgresRepository) getNextShortID(ctx context.Context) (int, error) {
	count := 0

	rows, err := repository.db.QueryContext(ctx, `SELECT COUNT(*) FROM "links"`)
	if err != nil {
		return 1, err
	}

	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 1, err
		}

		count = count + 1
	}

	err = rows.Err()
	if err != nil {
		return 1, err
	}

	return count, nil
}
