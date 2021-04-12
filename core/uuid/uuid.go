package uuid

import "github.com/google/uuid"

// return a random V4 uuid string
func New() string {
	return uuid.New().String()
}
