package db

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupDB() (*sql.DB, error) {
	dsn := ""
	db, err := sql.Open()
}
