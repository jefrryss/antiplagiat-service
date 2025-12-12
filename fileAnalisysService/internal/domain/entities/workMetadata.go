package entities

import "github.com/google/uuid"

type WorkMetadata struct {
	ID       uuid.UUID
	UserName string
	FileName string
}
