package manager

import (
	"context"
	"fileStoringService/internal/domain/entities"
	"fileStoringService/internal/domain/repository"
	"fmt"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type ManagerFileStorage struct {
	repo      repository.Repo
	s3Storage repository.Storage
}

func NewManagerFileStorage(repo repository.Repo, s3Storage repository.Storage) *ManagerFileStorage {
	return &ManagerFileStorage{repo: repo, s3Storage: s3Storage}
}

func (m *ManagerFileStorage) Save(ctx context.Context, userName, typeWork string, file multipart.File, header *multipart.FileHeader) (uuid.UUID, error) {
	if file == nil || header == nil {
		return uuid.Nil, fmt.Errorf("нет файла для загрузки")
	}

	workID := uuid.New()
	fileID := uuid.New()
	objectName := fileID.String()

	if err := m.s3Storage.SaveFile(objectName, file, header.Size); err != nil {
		return uuid.Nil, fmt.Errorf("ошибка сохранения файла в MinIO: %w", err)
	}

	f := entities.File{
		ID:          fileID,
		FileName:    header.Filename,
		ContentType: header.Header.Get("Content-Type"),
	}
	if f.ContentType == "" {
		f.ContentType = "application/octet-stream"
	}

	if err := m.repo.CreateFile(ctx, f); err != nil {
		_ = m.s3Storage.DeleteFile(objectName)
		return uuid.Nil, fmt.Errorf("ошибка сохранения файла в БД: %w", err)
	}

	w := entities.Work{
		ID:        workID,
		UserName:  userName,
		CreatedAt: time.Now(),
		TypeWork:  typeWork,
		File:      f,
	}

	if err := m.repo.CreateWork(ctx, w); err != nil {
		_ = m.s3Storage.DeleteFile(objectName)
		return uuid.Nil, fmt.Errorf("ошибка создания работы в БД: %w", err)
	}

	return workID, nil
}

func (m *ManagerFileStorage) GetWorksInfoByType(ctx context.Context, typeWork string) ([]entities.Work, error) {
	works, err := m.repo.GetWorksByType(ctx, typeWork)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения работ по типу %s: %w", typeWork, err)
	}

	if len(works) == 0 {
		return nil, fmt.Errorf("файлы не найдены для типа работы %s", typeWork)
	}

	return works, nil
}

func (m *ManagerFileStorage) GetWorkWithFile(ctx context.Context, workID string) (entities.Work, io.ReadCloser, error) {
	work, err := m.repo.GetWork(ctx, workID)
	if err != nil {
		return entities.Work{}, nil, fmt.Errorf("работа не найдена: %w", err)
	}

	fileReader, err := m.s3Storage.GetFile(work.File.ID.String())
	if err != nil {
		return entities.Work{}, nil, fmt.Errorf("не удалось получить файл из MinIO: %w", err)
	}

	return work, fileReader, nil
}
