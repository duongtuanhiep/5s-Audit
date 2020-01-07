package model

import "github.com/duongtuanhiep/fliqaudit/model/enum"

// CheckItemSend model, it define the type of checkitem object that will be sent to database
type CheckItemSend struct {
	QuestionID      int             `json:"questionID"`
	CheckItemAnswer enum.ItemAnswer `json:"checkItemAnswer"`
}
