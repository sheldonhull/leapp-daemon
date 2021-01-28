package services

import (
	"fmt"
)

func Home(name string) (string, error) {
	return fmt.Sprintf("hello %s", name), nil
}
