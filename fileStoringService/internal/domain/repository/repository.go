package repository

import (
	"context"
	"fileStoringService/internal/domain/entities"
)

type Repo interface {
	CreateFile(ctx context.Context, file entities.File) error
	CreateWork(ctx context.Context, work entities.Work) error
	GetWork(ctx context.Context, id string) (entities.Work, error)
	GetWorksByType(ctx context.Context, typeWork string) ([]entities.Work, error)
	Close()
}
