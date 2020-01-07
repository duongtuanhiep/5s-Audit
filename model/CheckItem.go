package model

import "github.com/duongtuanhiep/fliqaudit/model/enum"

// CheckItem model, it defines a single question and its answer. send back to client
type CheckItem struct {
	CheckItemID          int             `json:"checkItemID"`
	QuestionID           int             `json:"questionID"`
	CheckItem            string          `json:"checkItem"`
	CheckItemDescription string          `json:"checkItemDescription"`
	CheckItemAnswer      enum.ItemAnswer `json:"checkItemAnswer"`
	AuditPhase           enum.AuditPhase `json:"auditPhase"`
	Status               enum.Status     `json:"status"`
}
