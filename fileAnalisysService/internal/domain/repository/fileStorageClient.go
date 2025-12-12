package repository

import (
	"fileAnalisysService/internal/domain/entities"

	"github.com/google/uuid"
)

type FileStorageClient interface {
	GetWorkMetadataByType(typeWork string) ([]entities.WorkMetadata, error)
	GetWorkFile(workID uuid.UUID) ([]byte, error)
}
