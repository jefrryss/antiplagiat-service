package entities

import (
	"time"

	"github.com/google/uuid"
)

type FileData struct {
	FileID    uuid.UUID
	UserID    uuid.UUID
	UserName  string
	CreatedAt time.Time
}
