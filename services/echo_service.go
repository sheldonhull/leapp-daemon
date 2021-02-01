package services

import (
	"fmt"
)

func Echo(text string) (string, error) {
	return fmt.Sprintf("%s", text), nil
}
