package environment

import (
	"os"
	"runtime"
	"testing"
)

func TestGetEnvironmentVariable(t *testing.T) {
	variableName := "testVariable"
	variableValue := "check123"
	defer func(key string) {
		err := os.Unsetenv(key)
		if err != nil {
			t.Fatalf("Unable to unset the environment variable used for the test")
		}
	}(variableName)

	env := Environment{}
	if env.GetEnvironmentVariable(variableName) != "" {
		t.Fatalf("Unexpected variableValue")
	}

	err := os.Setenv(variableName, variableValue)
	if err != nil {
		t.Fatalf("Unable to set environment variable for the test")
	}

	if env.GetEnvironmentVariable(variableName) != variableValue {
		t.Fatalf("Unexpected variableValue")
	}
}

func TestIsCommandAvailable_available(t *testing.T) {
	env := Environment{}

	if runtime.GOOS == "windows" {
		if !env.IsCommandAvailable("cmd.exe") {
			t.Fatalf("cmd.exe should exist on Windows")
		}
	} else {
		if !env.IsCommandAvailable("ls") {
			t.Fatalf("ls should exist on Unix-Like system")
		}
	}
}

func TestIsCommandAvailable_unavailable(t *testing.T) {
	env := Environment{}

	if env.IsCommandAvailable("unavailable_command_23465236") {
		t.Fatalf("command should not exist")
	}
}

func TestGetOs(t *testing.T) {
	env := Environment{}

	if env.GetOs() == "" {
		t.Fatalf("Unable to define the current OS")
	}

	if env.GetOs() != runtime.GOOS {
		t.Fatalf("Unexpected OS returned")
	}
}

func TestIsWindows(t *testing.T) {
	env := Environment{}

	if runtime.GOOS == "windows" && !env.IsWindows() {
		t.Fatalf("Environment is Windows")
	} else if runtime.GOOS != "windows" && env.IsWindows() {
		t.Fatalf("Environment is not Windows")
	}
}
