package antiplagiat

import (
	"errors"
	"fileAnalisysService/internal/domain/antiplagiat"
)

type BitwiseEngine struct{}

func NewBitwiseEngine() antiplagiat.AntiPlagiarismEngine {
	return &BitwiseEngine{}
}

func (e *BitwiseEngine) Compare(a, b []byte) (float64, error) {
	if a == nil || b == nil {
		return 0, errors.New("входные данные не могут быть nil")
	}

	minLen := len(a)
	if len(b) < minLen {
		minLen = len(b)
	}
	if minLen == 0 {
		return 0, nil
	}

	same := 0
	for i := 0; i < minLen; i++ {
		if a[i] == b[i] {
			same++
		}
	}

	maxLen := len(a)
	if len(b) > maxLen {
		maxLen = len(b)
	}

	return float64(same) / float64(maxLen), nil
}
