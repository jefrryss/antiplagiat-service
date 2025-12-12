package entities

import "github.com/google/uuid"

type CompareResult struct {
	WorkA        uuid.UUID
	WorkB        uuid.UUID
	NameUserA    string
	NameUserB    string
	NameFileA    string
	NameFileB    string
	Similarity   float64
	IsPlagiarism bool
}
