package response

import "errors"

type ResponseError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *ResponseError) Error() string {
	return e.Msg
}

func IsResponseError(err error) (*ResponseError, bool) {
	var e *ResponseError
	if errors.As(err, &e) {
		return e, true
	}
	return nil, false
}

func NewResponseError(code int, msg string) *ResponseError {
	return &ResponseError{Code: code, Msg: msg}
}
