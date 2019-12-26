package models

type Result struct {
	ErrCode int         `json:"errcode"`
	ErrMsg  string      `json:"errmsg,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func NewDefaultResult() Result {
	return Result{}
}
