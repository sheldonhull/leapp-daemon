package logging

import (
	"os/user"
	"testing"
)

func TestGetHomeDir(t *testing.T) {
	usr, _ := user.Current()
	val, _ := GetHomeDir()
	if usr.HomeDir != val {
		t.Errorf("got %s, want %s", usr.HomeDir, val)
	}
}