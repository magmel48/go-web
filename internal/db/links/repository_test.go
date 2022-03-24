package links

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"regexp"
	"testing"
)

func TestPostgresRepository_Create(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx         context.Context
		shortID     string
		originalURL string
	}

	id := 1
	shortID := "1"
	originalURL := "https://google.com"

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectQuery(
		regexp.QuoteMeta(`
			INSERT INTO "links" ("short_id", "original_url") VALUES ($1, $2)
			ON CONFLICT ("original_url") DO UPDATE SET "original_url" = "links"."original_url"
			RETURNING "id", "short_id"
		`))
	e.WillReturnRows(sqlmock.NewRows([]string{"id", "short_id"}).AddRow(id, shortID))
	e.WillReturnError(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Link
		wantErr bool
	}{
		{
			name:    "should execute proper query",
			fields:  fields{db: db},
			args:    args{ctx: context.TODO(), shortID: shortID, originalURL: originalURL},
			want:    &Link{ID: id, ShortID: shortID},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgresRepository{
				db: tt.fields.db,
			}
			got, err := repository.Create(tt.args.ctx, tt.args.shortID, tt.args.originalURL)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostgresRepository_FindByShortID(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx     context.Context
		shortID string
	}

	id := 1
	shortID := "2"
	originalURL := "https://google.com"
	isDeleted := false

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectQuery(
		regexp.QuoteMeta(`SELECT "id", "short_id", "original_url", "is_deleted" FROM "links" WHERE "short_id" = $1 LIMIT 1`))
	e.WillReturnRows(sqlmock.NewRows([]string{"id", "short_id", "original_url", "is_deleted"}).AddRow(
		id, shortID, originalURL, isDeleted))
	e.WillReturnError(nil)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Link
		wantErr bool
	}{
		{
			name:    "should execute proper query",
			args:    args{ctx: context.TODO(), shortID: shortID},
			fields:  fields{db: db},
			want:    &Link{ShortID: shortID, ID: id, OriginalURL: originalURL, IsDeleted: isDeleted},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgresRepository{
				db: tt.fields.db,
			}
			got, err := repository.FindByShortID(tt.args.ctx, tt.args.shortID)

			if (err != nil) != tt.wantErr {
				t.Errorf("FindByShortID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FindByShortID() got = %v, want %v", got, tt.want)
			}
		})
	}
}
