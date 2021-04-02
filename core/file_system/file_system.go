package file_system

import (
	"leapp_daemon/custom_error"
	"os"
	"os/user"
)

func DoesFileExist(path string) bool {
	_, err := os.Stat(path)
	doesFileNotExists := os.IsNotExist(err)
	if doesFileNotExists { return false } else { return true }
}

func GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", custom_error.NewInternalServerError(err)
	}
	return usr.HomeDir, nil
}
