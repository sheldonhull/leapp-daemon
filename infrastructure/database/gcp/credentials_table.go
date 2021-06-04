package gcp

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"leapp_daemon/infrastructure/http/http_error"
)

var TableName = "credentials"

type CredentialsTable struct {
	accountId string
	value     string
}

func (table *CredentialsTable) InsertQuery(database *sql.DB, jsonObject string) error {
	insertQuery := fmt.Sprintf("INSERT INTO %v (account_id, value) VALUES(?, ?)", TableName)
	statement, err := database.Prepare(insertQuery)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	var cred CredentialsTable
	json.Unmarshal([]byte(jsonObject), &cred)

	_, err = statement.Exec(cred.accountId, cred.value)
	if err != nil {
		return http_error.NewBadRequestError(err)
	}

	return nil
}

func (table *CredentialsTable) ParseSqlLiteRows(rows *sql.Rows) ([]string, error) {
	defer rows.Close()

	var results []string
	for rows.Next() {
		var accountId string
		var value string
		rows.Scan(&accountId, &value)

		serializedRow, err := json.Marshal(CredentialsTable{accountId: accountId, value: value})
		if err != nil {
			return nil, http_error.NewInternalServerError(err)
		}

		results = append(results, string(serializedRow))
	}

	return results, nil
}
