package database

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"leapp_daemon/infrastructure/http/http_error"
)

type SqlLiteRowMapper interface {
	ParseSqlLiteRows(rows *sql.Rows) ([]string, error)
	InsertQuery(*sql.DB, string) error
}

type SqlLite struct {
	sqlLiteRowMapper SqlLiteRowMapper
}

func (sqlLite *SqlLite) insertRow(sqlFilePath string, jsonObject string) error {
	database, err := sql.Open("sqlite3", sqlFilePath)
	if err != nil {
		return http_error.NewNotFoundError(err)
	}

	defer database.Close()
	return sqlLite.sqlLiteRowMapper.InsertQuery(database, jsonObject)
}

func (sqlLite *SqlLite) executeQuery(slqLiteFilePath string, textQuery string) ([]string, error) {
	database, err := sql.Open("sqlite3", slqLiteFilePath)
	if err != nil {
		return nil, http_error.NewNotFoundError(err)
	}

	defer database.Close()
	rows, err := database.Query(textQuery)
	if err != nil {
		return nil, http_error.NewInternalServerError(err)
	}

	return sqlLite.sqlLiteRowMapper.ParseSqlLiteRows(rows)
}
