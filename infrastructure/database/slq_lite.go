package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"leapp_daemon/infrastructure/http/http_error"
)

type SqlLite struct {
}

func (sqlLite *SqlLite) GetDatabase(sqlFilePath string) (*sql.DB, error) {
	database, err := sql.Open("sqlite3", sqlFilePath)
	if err != nil {
		return nil, http_error.NewNotFoundError(err)
	}

	return database, nil
}
