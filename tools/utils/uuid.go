package utils

import (
	"github.com/google/uuid"
)

func Generate() string {
	return uuid.New().String()
}
func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}


func Parse(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}

func MustParse(id string) uuid.UUID {
	return uuid.MustParse(id)
}
