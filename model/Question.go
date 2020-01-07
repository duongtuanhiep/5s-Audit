package model

import "github.com/duongtuanhiep/fliqaudit/model/enum"

// Question is used to returned as a list to the client
type Question struct {
	QuestionID           int             `json:"questionID"`
	CheckItem            string          `json:"checkItem"`
	CheckItemDescription string          `json:"checkItemDescription"`
	AuditPhase           enum.AuditPhase `json:"auditPhase"`
	Status               enum.Status     `json:"status"`
}
