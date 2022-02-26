package userlinks

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/magmel48/go-web/internal/auth"
	"reflect"
	"regexp"
	"testing"
)

func TestNewPostgresRepository(t *testing.T) {
	type args struct {
		db *sql.DB
	}

	tests := []struct {
		name string
		args args
		want *PostgresRepository
	}{
		{
			name: "returns new instance of PostgresRepository",
			want: &PostgresRepository{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostgresRepository(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostgresRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
