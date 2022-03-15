package rest

import (
	"encoding/json"
	"errors"
	"net/http"
)

type RespBody struct {
	StatusCode int         `json:"-"`
	OK         bool        `json:"ok"`
	Data       interface{} `json:"data,omitempty"`
	ErrCode    string      `json:"err,omitempty"`
	Message    string      `json:"msg,omitempty"`
}

func NewSuccessResp(data interface{}) RespBody {
	return RespBody{
		StatusCode: http.StatusOK,
		OK:         true,
		Data:       data,
	}
}

func NewErrorResp(err error) RespBody {
	var e *Error
	if !errors.As(err, &e) {
		e = NewErrInternalError(err)
	}
	return RespBody{
		StatusCode: e.StatusCode,
		OK:         false,
		ErrCode:    e.ErrCode,
		Message:    e.Message,
	}
}

func WriteRespBody(w http.ResponseWriter, resp RespBody) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)

	b, _ := json.Marshal(resp)
	w.Write(b)
}
