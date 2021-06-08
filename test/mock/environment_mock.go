package mock

import (
	"fmt"
)

type EnvironmentMock struct {
	calls                 []string
	ExpIsCommandAvailable bool
	ExpIsWindows          bool
}

func NewEnvironmentMock() EnvironmentMock {
	return EnvironmentMock{calls: []string{}}
}

func (env *EnvironmentMock) GetCalls() []string {
	return env.calls
}

func (env *EnvironmentMock) GetEnvironmentVariable(variable string) string {
	env.calls = append(env.calls, fmt.Sprintf("GetEnvironmentVariable(%v)", variable))
	return "c:/appdata"
}

func (env *EnvironmentMock) IsCommandAvailable(command string) bool {
	env.calls = append(env.calls, fmt.Sprintf("IsCommandAvailable(%v)", command))
	return env.ExpIsCommandAvailable
}

func (env *EnvironmentMock) IsWindows() bool {
	env.calls = append(env.calls, "IsWindows()")
	return env.ExpIsWindows
}