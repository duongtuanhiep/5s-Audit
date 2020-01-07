package model

import (
	// Well fuck you go
	_ "encoding/json"
)

// AuditStat : Phase average score of all phase and audit id
type AuditStat struct {
	Key         int            `json:"key"`
	PhaseScores []PhaseAverage `json:"phaseScores"`
}
