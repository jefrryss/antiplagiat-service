package entities

import "github.com/google/uuid"

type File struct {
	ID          uuid.UUID
	FileName    string
	ContentType string
}
