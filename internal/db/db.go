package db

import (
	"context"
	"database/sql"
	"errors"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magmel48/go-web/internal/config"
	"log"
	"time"
)

// DB is common interface that is designed to wrap a work with SQL database.
//go:generate mockery --name=DB
type DB interface {
	Instance() *sql.DB
	Connect() error
	CheckConnection(ctx context.Context) bool
	CreateSchema() error
}

// SQLDB is implementation of abstract DB.
type SQLDB struct {
	instance *sql.DB
}

// Instance returns current *sql.DB instance that manages connections to SQL database.
func (db *SQLDB) Instance() *sql.DB {
	return db.instance
}

// Connect opens connection to SQL database.
func (db *SQLDB) Connect() error {
	if config.DatabaseDSN != "" {
		var err error

		db.instance, err = sql.Open("pgx", config.DatabaseDSN)
		return err
	}

	return errors.New("no database config specified")
}

// CheckConnection checks if a new connection to database can be potentially opened
// and everything related to DB was configured properly previously.
func (db *SQLDB) CheckConnection(ctx context.Context) bool {
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	err := db.instance.PingContext(ctx)

	if err != nil {
		log.Println("db connection error", err)
	}

	return err == nil
}
