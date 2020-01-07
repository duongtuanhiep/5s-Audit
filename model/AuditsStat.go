package model

import (
	// Well fuck you go
	_ "encoding/json"
)

// AuditStat : Phase average score of all phase and audit id
type AuditsStat struct {
	PhaseScores []PhaseAverage `json:"phaseScores"`
}
