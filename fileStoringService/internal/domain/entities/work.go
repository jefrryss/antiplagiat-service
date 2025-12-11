package entities

import (
	"time"

	"github.com/google/uuid"
)

type Work struct {
	ID        uuid.UUID
	UserName  string
	CreatedAt time.Time
	File      File
	TypeWork  string
}
