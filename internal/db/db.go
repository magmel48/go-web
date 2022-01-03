package db

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/magmel48/go-web/internal/config"
	"log"
	"time"
)

var DB *sql.DB

func Connect() {
	if config.DatabaseDSN != "" {
		var err error

		DB, err = sql.Open("pgx", config.DatabaseDSN)
		if err != nil {
			panic(err)
		}
	}
}

func CheckConnection() bool {
	ctx, cancel := context.WithTimeout(context.TODO(), 2*time.Second)
	defer cancel()
	err := DB.PingContext(ctx)

	if err != nil {
		log.Println("db connection error", err)
	}

	return err == nil
}
