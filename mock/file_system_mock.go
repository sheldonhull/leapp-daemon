package mock

import (
	"errors"
	"fmt"
)

type FileSystemMock struct {
	calls                []string
	ExpErrorOnGetHomeDir bool
	ExpDoesFileExist     bool
}

func (fileSystem *FileSystemMock) GetCalls() []string {
	return fileSystem.calls
}

func NewFileSystemMock() FileSystemMock {
	return FileSystemMock{calls: []string{}}
}

func (fileSystem *FileSystemMock) DoesFileExist(path string) bool {
	fileSystem.calls = append(fileSystem.calls, fmt.Sprintf("DoesFileExist(%v)", path))
	return fileSystem.ExpDoesFileExist
}

func (fileSystem *FileSystemMock) GetHomeDir() (string, error) {
	if fileSystem.ExpErrorOnGetHomeDir {
		return "", errors.New("error")
	}
	fileSystem.calls = append(fileSystem.calls, "GetHomeDir()")
	return "/user/home", nil
}
