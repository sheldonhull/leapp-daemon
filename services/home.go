package services

import (
	"github.com/pkg/errors"
)

func Home(name string) (string, error) {
	return "", errors.New("fake error")
	//return fmt.Sprintf("hello %s", name), nil
}
