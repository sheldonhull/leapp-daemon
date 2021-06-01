package environment

import (
	"os"
	"os/exec"
	"runtime"
)

type Environment struct {
}

func (env *Environment) GetEnvironmentVariable(variable string) string {
	return os.Getenv(variable)
}

func (env *Environment) IsCommandAvailable(command string) bool {
	_, err := exec.LookPath(command)
	return err == nil
}

func (env *Environment) GetOs() string {
	return runtime.GOOS
}

func (env *Environment) IsWindows() bool {
	return env.GetOs() == "windows"
}
