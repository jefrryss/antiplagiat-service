package entities

import (
	"time"

	"github.com/google/uuid"
)

type AnalysisReport struct {
	ReportID  uuid.UUID
	TypeWork  string
	CreatedAt time.Time

	Results []CompareResult
}
