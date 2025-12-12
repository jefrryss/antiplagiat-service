package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type ReportManagerImpl struct {
	analysisServiceURL string
	client             *http.Client
}

func NewReportManager(analysisServiceURL string) *ReportManagerImpl {
	return &ReportManagerImpl{
		analysisServiceURL: analysisServiceURL,
		client:             &http.Client{},
	}
}

func (r *ReportManagerImpl) CreateReport(typeWork, workID string) (map[string]interface{}, error) {
	reqBody := map[string]string{"work_id": workID}
	reqJSON, _ := json.Marshal(reqBody)

	resp, err := r.client.Post(r.analysisServiceURL+"/works/"+typeWork+"/reports", "application/json", bytes.NewReader(reqJSON))
	if err != nil {
		return nil, errors.New("analysis service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, errors.New("failed to create report")
	}

	var reportResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&reportResp); err != nil {
		return nil, err
	}

	return reportResp, nil
}

func (r *ReportManagerImpl) GetLatestReport(typeWork string) (map[string]interface{}, error) {
	resp, err := r.client.Get(r.analysisServiceURL + "/works/" + typeWork + "/reports/last")
	if err != nil {
		return nil, errors.New("analysis service unavailable")
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, errors.New("report not found")
	}
	if resp.StatusCode != 200 {
		return nil, errors.New("error from analysis service")
	}

	var report map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, err
	}

	return report, nil
}
