package model

import (

	// Well fuck you go
	_ "encoding/json"
	"time"
)

// Audit model, used this to jsonify and send back to client
type Audit struct {
	Key            int         `json:"key"`
	DepartmentName string      `json:"departmentName"`
	CompanyName    string      `json:"companyName"`
	AuditorName    string      `json:"auditorName"`
	AssessedDate   time.Time   `json:"assessedDate"`
	CheckItems     []CheckItem `json:"checkItems"`
	// CheckItems     []CheckItem `json:"checkItems"`
}
