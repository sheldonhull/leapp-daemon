package gcp

import (
  "database/sql"
  "io/ioutil"
  "leapp_daemon/test"
  "os"
  "path"
  "reflect"
  "testing"
)

type CredentialsRow struct {
  accountId string
  value     string
}

var (
  tempDbPath string
)

func setup() {
  tempDbPath = path.Join(os.TempDir(), "sql_lite_test.db")
}

func expectCredentials(t *testing.T, expectedRows []CredentialsRow) {
  database, _ := sql.Open("sqlite3", tempDbPath)
  defer database.Close()

  row, _ := database.Query("SELECT * FROM credentials ORDER BY account_id ASC")
  defer row.Close()

  foundRows := []CredentialsRow{}
  var dbAccountId string
  var dbValue string
  for row.Next() {
    row.Scan(&dbAccountId, &dbValue)
    foundRows = append(foundRows, CredentialsRow{accountId: dbAccountId, value: dbValue})
  }

  if !reflect.DeepEqual(foundRows, expectedRows) {
    t.Errorf("expected rows:\n%v", expectedRows)
  }
}

func TestWriteCredentials(t *testing.T) {
  setup()

  accountId := "account_id@domain.com"
  anotherAccountId := "another_account@domain.com"
  value := "jsonCredentials"
  anotherValue := "anotherJsonCredentials"

  createdFile, _ := os.Create(tempDbPath)
  createdFile.Close()

  table := CredentialsTable{}
  table.WriteCredentials(tempDbPath, accountId, value)
  expectCredentials(t, []CredentialsRow{{accountId: accountId, value: value}})

  table.WriteCredentials(tempDbPath, anotherAccountId, value)
  expectCredentials(t, []CredentialsRow{{accountId: accountId, value: value}, {accountId: anotherAccountId, value: value}})

  table.WriteCredentials(tempDbPath, accountId, anotherValue)
  expectCredentials(t, []CredentialsRow{{accountId: accountId, value: anotherValue}, {accountId: anotherAccountId, value: value}})
}

func TestWriteCredentials_CorruptedDbFile(t *testing.T) {
  setup()
  os.Remove(tempDbPath)
  ioutil.WriteFile(tempDbPath, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0644)

  accountId := "account_id@domain.com"
  value := "jsonCredentials"

  table := CredentialsTable{}
  err := table.WriteCredentials(tempDbPath, accountId, value)
  test.ExpectHttpError(t, err, 500, "file is not a database")
}

func TestRemoveCredentials(t *testing.T) {
  setup()

  accountId := "account_id@domain.com"
  anotherAccountId := "another_account@domain.com"
  value := "jsonCredentials"
  anotherValue := "anotherJsonCredentials"

  createdFile, _ := os.Create(tempDbPath)
  createdFile.Close()

  table := CredentialsTable{}
  table.WriteCredentials(tempDbPath, accountId, value)
  table.WriteCredentials(tempDbPath, anotherAccountId, anotherValue)

  table.RemoveCredentials(tempDbPath, anotherAccountId)
  expectCredentials(t, []CredentialsRow{{accountId: accountId, value: value}})
}

func TestRemoveCredentials_CorruptedDbFile(t *testing.T) {
  setup()
  os.Remove(tempDbPath)
  ioutil.WriteFile(tempDbPath, []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 0644)

  table := CredentialsTable{}
  err := table.RemoveCredentials(tempDbPath, "account_id@domain.com")
  test.ExpectHttpError(t, err, 500, "file is not a database")
}