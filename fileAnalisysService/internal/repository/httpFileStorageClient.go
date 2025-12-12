package repository

import (
	"encoding/json"
	"fileAnalisysService/internal/domain/entities"
	"fileAnalisysService/internal/domain/repository"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
)

type HttpFileStorageClient struct {
	baseURL string
	client  *http.Client
}

func NewHttpFileStorageClient(baseURL string) repository.FileStorageClient {
	return &HttpFileStorageClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *HttpFileStorageClient) GetWorkMetadataByType(typeWork string) ([]entities.WorkMetadata, error) {
	url := fmt.Sprintf("%s/files/list/%s", c.baseURL, typeWork)

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса к FileStorage: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("FileStorage вернул %s: %s", resp.Status, string(body))
	}

	var works []struct {
		ID       string `json:"id"`
		UserName string `json:"userName"`
		File     struct {
			FileName string `json:"fileName"`
		} `json:"file"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&works); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	var metadata []entities.WorkMetadata
	for _, w := range works {
		uid, err := uuid.Parse(w.ID)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга UUID: %w", err)
		}
		metadata = append(metadata, entities.WorkMetadata{
			ID:       uid,
			UserName: w.UserName,
			FileName: w.File.FileName,
		})
	}

	return metadata, nil
}

func (c *HttpFileStorageClient) GetWorkFile(workID uuid.UUID) ([]byte, error) {
	url := fmt.Sprintf("%s/files/download/%s", c.baseURL, workID.String())

	resp, err := c.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка скачивания файла: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("FileStorage вернул %s: %s", resp.Status, string(body))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	return data, nil
}
