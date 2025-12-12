package repository

import (
	"bytes"
	"context"
	"fileAnalisysService/internal/domain/repository"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioReportStorage struct {
	client *minio.Client
	bucket string
}

func NewMinioReportStorage(endpoint, accessKey, secretKey, bucket string) repository.StorageReports {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Ошибка подключения к MinIO %v", err)
	}

	return &MinioReportStorage{client: client, bucket: bucket}
}

func (m *MinioReportStorage) SaveReport(objectName string, data []byte) error {
	if objectName == "" {
		return fmt.Errorf("имя объекта пустое")
	}

	ctx := context.Background()

	exists, err := m.client.BucketExists(ctx, m.bucket)
	if err != nil {
		return fmt.Errorf("ошибка проверки бакета %w", err)
	}

	if !exists {
		if err := m.client.MakeBucket(ctx, m.bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("не удалось создать бакет %s: %w", m.bucket, err)
		}
	}

	reader := bytes.NewReader(data)

	_, err = m.client.PutObject(ctx, m.bucket, objectName, reader, int64(len(data)), minio.PutObjectOptions{
		ContentType: "application/json",
	})
	if err != nil {
		return fmt.Errorf("ошибка загрузки отчёта в MinIO %w", err)
	}

	return nil
}

func (m *MinioReportStorage) GetReport(objectName string) (io.ReadCloser, error) {
	if objectName == "" {
		return nil, fmt.Errorf("имя объекта пустое")
	}

	obj, err := m.client.GetObject(context.Background(), m.bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения отчёта из MinIO %w", err)
	}

	if _, err = obj.Stat(); err != nil {
		return nil, fmt.Errorf("отчёт не найден %w", err)
	}

	return obj, nil
}

func (m *MinioReportStorage) DeleteReport(objectName string) error {
	if objectName == "" {
		return fmt.Errorf("имя объекта пустое")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := m.client.RemoveObject(ctx, m.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("ошибка удаления отчёта %w", err)
	}

	return nil
}
