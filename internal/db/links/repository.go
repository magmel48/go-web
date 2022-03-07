package links

import (
	"context"
	"database/sql"
	"log"
	"strconv"
)

// PostgresRepository is implementation of abstract Repository.
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository returns new PostgresRepository for working with links.
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

// Create creates new shorter link by specified originalURL.
func (repository *PostgresRepository) Create(ctx context.Context, shortID string, originalURL string) (*Link, error) {
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

		return nil, err
	}

	// shortID != link.ShortID if short_id`s are not the same
	var err error
	if shortID != link.ShortID {
		err = ErrConflict
	}

	return &link, err
}

// CreateBatch creates many shorter links by specified originalURLs.
func (repository *PostgresRepository) CreateBatch(ctx context.Context, originalURLs []string) ([]Link, error) {
	result := make([]Link, len(originalURLs))
	linksCount, _ := repository.getNextShortID(ctx)

	tx, err := repository.db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		err := tx.Rollback()
		log.Println("CreateBatch tx rollback error", err)
	}()

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

// FindByShortID finds originalURL and related info by specified short link identifier.
func (repository *PostgresRepository) FindByShortID(ctx context.Context, shortID string) (*Link, error) {
	rows, err := repository.db.QueryContext(
		ctx, `SELECT "id", "short_id", "original_url", "is_deleted" FROM "links" WHERE "short_id" = $1 LIMIT 1`, shortID)
	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("FindByShortID close rows error", err)
		}
	}()

	link := Link{}
	if rows.Next() {
		if err := rows.Scan(&link.ID, &link.ShortID, &link.OriginalURL, &link.IsDeleted); err != nil {
			return nil, err
		}
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	if link.ID != 0 {
		return &link, nil
	}

	return nil, nil
}

// getNextShortID returns next shortID that can be used for a new shorter link.
func (repository *PostgresRepository) getNextShortID(ctx context.Context) (int, error) {
	count := 0

	rows, err := repository.db.QueryContext(ctx, `SELECT COUNT(*) FROM "links"`)
	if err != nil {
		return 1, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Println("getNextShortID close rows error", err)
		}
	}()

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
