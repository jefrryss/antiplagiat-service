package manager

import (
	"apiGateway/internal/domain/manager"
	"errors"
	"io"
)

type PlagiarismService struct {
	fileMgr      manager.FileManager
	reportMgr    manager.ReportManager
	wordCloudMgr manager.WordCloudManager
}

func NewPlagiarismService(f manager.FileManager, r manager.ReportManager, w manager.WordCloudManager) *PlagiarismService {
	return &PlagiarismService{
		fileMgr:      f,
		reportMgr:    r,
		wordCloudMgr: w,
	}
}
func (s *PlagiarismService) GetLatestReport(typeWork string) (map[string]interface{}, error) {
	return s.reportMgr.GetLatestReport(typeWork)
}
func (s *PlagiarismService) UploadFileAndGetWordCloud(userName, typeWork, fileName string, fileData io.Reader) (string, error) {
	workID, err := s.fileMgr.UploadFile(userName, typeWork, fileName, fileData)
	if err != nil {
		return "", err
	}

	_, err = s.reportMgr.CreateReport(typeWork, workID)
	if err != nil {
		return "", err
	}

	if seeker, ok := fileData.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	}

	fileBytes, err := io.ReadAll(fileData)
	if err != nil {
		return "", err
	}

	text := string(fileBytes)
	if text == "" {
		return "", errors.New("file is empty, cannot generate word cloud")
	}

	cloudURL := s.wordCloudMgr.GenerateWordCloud(text)
	return cloudURL, nil
}
