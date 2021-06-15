package gcp

import (
  "database/sql"
  "fmt"
  _ "github.com/mattn/go-sqlite3"
  "leapp_daemon/infrastructure/http/http_error"
)

var accessTokensTableName = "access_tokens"

type AccessTokensTable struct {
}

func (table *AccessTokensTable) RemoveAccessToken(sqlFilePath string, accountId string) error {
	database, err := getSqliteDatabase(sqlFilePath)
	if err != nil {
		return err
	}
	defer database.Close()

	err = table.deleteQuery(database, accountId)
	if err != nil {
		return err
	}

	return nil
}

func (table *AccessTokensTable) deleteQuery(database *sql.DB, accountId string) error {
	deleteQuery := fmt.Sprintf("DELETE FROM %v WHERE account_id=?", accessTokensTableName)
	statement, err := database.Prepare(deleteQuery)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	_, err = statement.Exec(accountId)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	return nil
}
