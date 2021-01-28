package services

import (
	"fmt"
)

func HomeService(name string) (string, error) {
	return fmt.Sprintf("hello %s", name), nil
}
