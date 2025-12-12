package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
)

type FileManagerImpl struct {
	fileServiceURL string
	client         *http.Client
}

func NewFileManager(fileServiceURL string) *FileManagerImpl {
	return &FileManagerImpl{
		fileServiceURL: fileServiceURL,
		client:         &http.Client{},
	}
}

func (f *FileManagerImpl) UploadFile(userName, typeWork, fileName string, fileData io.Reader) (string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", fileName)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(part, fileData)
	if err != nil {
		return "", err
	}

	writer.WriteField("userName", userName)
	writer.WriteField("typeWork", typeWork)
	writer.Close()

	req, err := http.NewRequest("POST", f.fileServiceURL+"/upload", body)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := f.client.Do(req)
	if err != nil {
		return "", errors.New("file service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("failed to upload file")
	}

	var respData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&respData); err != nil {
		return "", err
	}

	workID, ok := respData["work_id"].(string)
	if !ok {
		return "", errors.New("invalid response from file service")
	}

	return workID, nil
}
