package mock

import (
  "errors"
  "fmt"
)

type FileSystemMock struct {
  errorOnGetHomeDir bool
  calls             []string
}

func (fileSystem *FileSystemMock) GetCalls() []string {
  return fileSystem.calls
}

func NewFileSystemMock(errorOnGetHomeDir bool) FileSystemMock {
  return FileSystemMock{errorOnGetHomeDir: errorOnGetHomeDir, calls: []string{}}
}

func (fileSystem *FileSystemMock) DoesFileExist(path string) bool {
  fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("DoesFileExist(\"%v\")", path))
  return false
}

func (fileSystem *FileSystemMock) GetHomeDir() (string, error) {
  if fileSystem.errorOnGetHomeDir {
    return "", errors.New("error")
  }
  fileSystem.calls = append(fileSystem.calls, "GetHomeDir()")
  return "/user/home", nil
}
