package benchmarking

import (
	"database/sql"
	"github.com/magmel48/go-web/internal/config"
	"github.com/magmel48/go-web/internal/db"
	"os"
)

func DBConnect() (*sql.DB, error) {
	config.DatabaseDSN = os.Getenv("DATABASE_DSN")
	database := db.SQLDB{}
	err := database.Connect()

	return database.Instance(), err
}
