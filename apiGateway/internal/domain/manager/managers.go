package manager

import "io"

type FileManager interface {
	UploadFile(userName, typeWork, fileName string, fileData io.Reader) (workID string, err error)
}

type ReportManager interface {
	CreateReport(typeWork, workID string) (report map[string]interface{}, err error)
	GetLatestReport(typeWork string) (report map[string]interface{}, err error)
}

type WordCloudManager interface {
	GenerateWordCloud(text string) string
}
