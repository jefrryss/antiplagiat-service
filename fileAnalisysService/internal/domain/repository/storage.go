package repository

import "io"

type StorageReports interface {
	SaveReport(objectName string, data []byte) error
	GetReport(objectName string) (io.ReadCloser, error)
	DeleteReport(objectName string) error
}
