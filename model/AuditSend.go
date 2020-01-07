package model

import (

	// Well fuck you go

	// Well fuck you go
	_ "encoding/json"
	"time"
)

// AuditSend model, used this to jsonify and send back to client
type AuditSend struct {
	DepartmentName string          `json:"departmentName"`
	CompanyName    string          `json:"companyName"`
	AuditorName    string          `json:"auditorName"`
	AssessedDate   time.Time       `json:"assessedDate"`
	CheckItemsSend []CheckItemSend `json:"checkItemsSend"`
	// CheckItems     []CheckItem `json:"checkItems"`
}
