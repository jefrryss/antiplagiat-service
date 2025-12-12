package repository

import (
	"context"
	"fmt"
	"io"
	"time"

	"fileStoringService/internal/domain/repository"

	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioStorage struct {
	Client *minio.Client
	bucket string
}

func NewMinioStorage(endpoint, accessKey, secretKey, bucket string) repository.Storage {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Не удалось подключиться к MinIO %v", err)
	}

	return &MinioStorage{
		Client: client,
		bucket: bucket,
	}
}

func (m *MinioStorage) SaveFile(id string, data io.Reader, size int64) error {
	if id == "" {
		return fmt.Errorf("имя объекта не может быть пустым")
	}
	if data == nil {
		return fmt.Errorf("данные файла пусты")
	}

	ctx := context.Background()

	exists, err := m.Client.BucketExists(ctx, m.bucket)
	if err != nil {
		return fmt.Errorf("ошибка проверки существования бакета %w", err)
	}

	if !exists {
		if err := m.Client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("не удалось создать бакет %s %w", m.bucket, err)
		}
	}

	if _, err := m.Client.PutObject(ctx, m.bucket, id, data, size, minio.PutObjectOptions{}); err != nil {
		return fmt.Errorf("ошибка загрузки файла в MinIO: %w", err)
	}

	return nil
}

func (m *MinioStorage) GetFile(objectName string) (io.ReadCloser, error) {
	if objectName == "" {
		return nil, fmt.Errorf("имя объекта не может быть пустым")
	}

	obj, err := m.Client.GetObject(context.Background(), m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("не удалось получить файл из MinIO %w", err)
	}

	_, err = obj.Stat()
	if err != nil {
		return nil, fmt.Errorf("файл не найден в MinIO %w", err)
	}

	return obj, nil
}

func (m *MinioStorage) DeleteFile(objectName string) error {
	if objectName == "" {
		return fmt.Errorf("имя объекта не может быть пустым")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := m.Client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{}); err != nil {
		return fmt.Errorf("не удалось удалить файл из MinIO %w", err)
	}

	return nil
}
