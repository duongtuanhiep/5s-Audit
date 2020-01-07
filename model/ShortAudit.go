package model

import "time"

// ShortAudit model. Send as a list of these item to client
type ShortAudit struct {
	Key          int       `json:"key"`
	AuditorName  string    `json:"auditorName"`
	AssessedDate time.Time `json:"assessedDate"`
	AverageScore float32   `json:"averageScore"`
}
