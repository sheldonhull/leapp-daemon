package mock

import (
	"errors"
	"fmt"
)

type GcpAccessTokensDbTableMock struct {
	calls                     []string
	ExpErrorOnExecDeleteQuery bool
}

func NewGcpAccessTokensDbTableMock() GcpAccessTokensDbTableMock {
	return GcpAccessTokensDbTableMock{calls: []string{}}
}

func (table *GcpAccessTokensDbTableMock) GetCalls() []string {
	return table.calls
}

func (table *GcpAccessTokensDbTableMock) RemoveAccessToken(sqlFilePath string, accountId string) error {
	if table.ExpErrorOnExecDeleteQuery {
		return errors.New("error executing delete query")
	}

	table.calls = append(table.calls, fmt.Sprintf("RemoveAccessToken(%v, %v)", sqlFilePath, accountId))
	return nil
}
