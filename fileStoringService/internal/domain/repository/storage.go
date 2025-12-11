package repository

import "io"

type Storage interface {
	SaveFile(objectName string, data io.Reader, size int64) error
	GetFile(objectName string) (io.ReadCloser, error)
	DeleteFile(objectName string) error
}
