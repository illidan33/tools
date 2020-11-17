package httptool

import (
	"encoding/json"
)

type Result struct {
	Code ResultCode      `json:"code"`
	Msg  string          `json:"msg"`
	Data json.RawMessage `json:"data"`
	Time int64           `json:"time"`
}

type ResultCode int

const (
	RESULT_SUCCESS ResultCode = 0
)
