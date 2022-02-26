package userlinks

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/magmel48/go-web/internal/auth"
	"github.com/magmel48/go-web/internal/db/links"
	"reflect"
	"regexp"
	"strconv"
	"testing"
)

func TestPostgresRepository_Create(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID auth.UserID
		linkID int
	}

	userID := "test_user_id"
	linkID := 52

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectExec(
		regexp.QuoteMeta(`INSERT INTO "user_links" ("user_id", "link_id") VALUES ($1, $2)`))
	e.WillReturnResult(sqlmock.NewResult(1, 1))
	e.WillReturnError(nil)
	e.WithArgs(userID, linkID)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "should execute proper query",
			fields:  fields{db: db},
			args:    args{ctx: context.TODO(), userID: &userID, linkID: linkID},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgresRepository{
				db: tt.fields.db,
			}

			if err := repository.Create(tt.args.ctx, tt.args.userID, tt.args.linkID); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostgresRepository_List(t *testing.T) {
	type fields struct {
		db *sql.DB
	}
	type args struct {
		ctx    context.Context
		userID auth.UserID
	}

	userID := "test_user_id"
	shortID := "2"
	originalURL := "https://google.com"

	db, sqlMock, _ := sqlmock.New()
	e := sqlMock.ExpectQuery(
		regexp.QuoteMeta(
			`SELECT l."short_id", l."original_url" FROM "user_links" AS ul JOIN "links" as l ON ul."link_id" = l."id" WHERE ul."user_id" = $1`))
	e.WillReturnRows(sqlmock.NewRows([]string{"short_id", "original_url"}).AddRow(shortID, originalURL))
	e.WillReturnError(nil)
	e.WithArgs(userID)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []UserLink
		wantErr bool
	}{
		{
			name:   "should execute proper query",
			fields: fields{db: db},
			args:   args{ctx: context.TODO(), userID: &userID},
			want:   []UserLink{{UserID: &userID, Link: links.Link{ShortID: shortID, OriginalURL: originalURL}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repository := &PostgresRepository{
				db: tt.fields.db,
			}
			got, err := repository.List(tt.args.ctx, tt.args.userID)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkPostgresRepository_DeleteLinks(b *testing.B) {
	db, _, _ := sqlmock.New()
	repository := PostgresRepository{db: db}

	userID := "test_user_id"
	shortIDs := make([]string, 10000)
	for i := 0; i < 10000; i++ {
		shortIDs[i] = strconv.Itoa(i + 1)
	}

	items := []DeleteQueryItem{
		{UserID: &userID, ShortIDs: shortIDs},
		{UserID: &userID, ShortIDs: shortIDs},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repository.DeleteLinks(context.TODO(), items)
	}
}
