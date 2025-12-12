package manager

import (
	"encoding/json"
	"fileAnalisysService/internal/domain/antiplagiat"
	"fileAnalisysService/internal/domain/entities"
	"fileAnalisysService/internal/domain/repository"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
)

type FileAnalysisManager struct {
	fsClient repository.FileStorageClient
	engine   antiplagiat.AntiPlagiarismEngine
	s3       repository.StorageReports
}

func NewFileAnalysisManager(fs repository.FileStorageClient, anti antiplagiat.AntiPlagiarismEngine, repo repository.StorageReports) *FileAnalysisManager {
	return &FileAnalysisManager{
		fsClient: fs,
		engine:   anti,
		s3:       repo,
	}
}
func (m *FileAnalysisManager) AnalyzeTypeWork(typeWork string) (*entities.AnalysisReport, error) {
	workMetadata, err := m.fsClient.GetWorkMetadataByType(typeWork)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения метаданных работ: %w", err)
	}

	report := &entities.AnalysisReport{
		ReportID:  uuid.New(),
		TypeWork:  typeWork,
		CreatedAt: time.Now(),
	}

	files := make(map[uuid.UUID][]byte)
	users := make(map[uuid.UUID]string)
	names := make(map[uuid.UUID]string)

	for _, meta := range workMetadata {
		data, err := m.fsClient.GetWorkFile(meta.ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка скачивания файла %s: %w", meta.ID, err)
		}
		files[meta.ID] = data
		users[meta.ID] = meta.UserName
		names[meta.ID] = meta.FileName
	}

	for i := 0; i < len(workMetadata); i++ {
		for j := i + 1; j < len(workMetadata); j++ {
			aID := workMetadata[i].ID
			bID := workMetadata[j].ID

			similarity, err := m.engine.Compare(files[aID], files[bID])
			if err != nil {
				return nil, fmt.Errorf("ошибка сравнения %s и %s: %w", aID, bID, err)
			}

			result := entities.CompareResult{
				WorkA:        aID,
				WorkB:        bID,
				NameUserA:    users[aID],
				NameUserB:    users[bID],
				NameFileA:    names[aID],
				NameFileB:    names[bID],
				Similarity:   similarity,
				IsPlagiarism: similarity > 0.8,
			}

			report.Results = append(report.Results, result)
		}
	}

	data, err := json.Marshal(report)
	if err != nil {
		return nil, fmt.Errorf("ошибка маршалинга отчета: %w", err)
	}

	err = m.s3.SaveReport(typeWork, data)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения отчета в S3: %w", err)
	}

	return report, nil
}

func (m *FileAnalysisManager) GetLastReportByTypeWork(typeWork string) (*entities.AnalysisReport, error) {
	obj, err := m.s3.GetReport(typeWork)
	if err != nil {
		return nil, fmt.Errorf("не удалось получить отчет из S3: %w", err)
	}
	defer obj.Close()

	data, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения отчета: %w", err)
	}

	var report entities.AnalysisReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("ошибка разбора JSON отчета: %w", err)
	}

	return &report, nil
}
