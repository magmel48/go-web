package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magmel48/go-web/internal/config"
	"log"
	"time"
)

//go:generate mockery --name=DB
type DB interface {
	Instance() *sql.DB
	Connect() error
	CheckConnection(ctx context.Context) bool
	CreateSchema() error
}

// SqlDB is implementation of abstract DB.
type SqlDB struct {
	instance *sql.DB
}

func (db *SqlDB) Instance() *sql.DB {
	return db.instance
}

func (db *SqlDB) Connect() error {
	if config.DatabaseDSN != "" {
		var err error

		db.instance, err = sql.Open("pgx", config.DatabaseDSN)
		return err
	}

	return nil
}

func (db *SqlDB) CheckConnection(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := db.instance.PingContext(ctx)

	if err != nil {
		log.Println("db connection error", err)
	}

	return err == nil
}
