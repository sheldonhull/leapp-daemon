package mock

import "fmt"

type EnvironmentMock struct {
  isWindows bool
  calls     []string
}

func NewEnvironmentMock(isWindows bool) EnvironmentMock {
  return EnvironmentMock{isWindows: isWindows, calls: []string{}}
}

func (env *EnvironmentMock) GetCalls() []string {
  return env.calls
}

func (env *EnvironmentMock) GetEnvironmentVariable(variable string) string {
  env.calls = append(env.calls, fmt.Sprintf("GetEnvironmentVariable(\"%v\")", variable))
  return "c:/appdata"
}

func (env *EnvironmentMock) IsWindows() bool {
  env.calls = append(env.calls, "IsWindows()")
  return env.isWindows
}
