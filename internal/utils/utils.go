package utils

import "github.com/google/uuid"

func NewUuidBytes() []byte {
	id := uuid.New()
	return id[:]
}
