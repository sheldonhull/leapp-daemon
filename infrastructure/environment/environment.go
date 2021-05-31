package environment

import (
  "os"
  "runtime"
)

type Environment struct {
}

func (env *Environment) GetEnvironmentVariable(variable string) string {
  return os.Getenv(variable)
}

func (env *Environment) GetOs() string {
  return runtime.GOOS
}

func (env *Environment) IsWindows() bool {
  return env.GetOs() == "windows"
}
