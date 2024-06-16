package postgres

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"library/internal/config"
)

const (
	booksTable   = "books"
	authorsTable = "authors"
)

func NewPostgresDB(cfg *config.DbServer) (*sql.DB, error) {

	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))

	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil

}
