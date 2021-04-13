package file_system

import (
  "io/ioutil"
  "os"
  "os/user"
)

type FileSystem struct {}

func(fileSystem *FileSystem) DoesFileExist(path string) bool {
	_, err := os.Stat(path)
	doesFileNotExists := os.IsNotExist(err)
	if doesFileNotExists {
	  return false
	} else {
	  return true
	}
}

func(fileSystem *FileSystem) GetHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func(fileSystem *FileSystem) ReadFile(path string) ([]byte, error) {
  encryptedText, err := ioutil.ReadFile(path)
  if err != nil {
    return nil, err
  }
  return encryptedText, nil
}

func(fileSystem *FileSystem) RemoveFile(path string) error {
  err := os.Remove(path)
  if err != nil {
    return err
  }
  return nil
}

func(fileSystem *FileSystem) WriteFile(path string, data []byte) error {
  err := ioutil.WriteFile(path, data, 0644)
  if err != nil {
    return err
  }
  return nil
}
