package model

import (
	"github.com/duongtuanhiep/fliqaudit/model/enum"

	// Well fuck you go
	_ "encoding/json"
)

// PhaseAverage : a single phase average including phase and average score of that phase
type PhaseAverage struct {
	AuditPhase enum.AuditPhase `json:"auditPhase"`
	Score      float32         `json:"score"`
	// CheckItems     []CheckItem `json:"checkItems"`
}
