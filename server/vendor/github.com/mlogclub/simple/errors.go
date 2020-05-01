package simple

import (
	"strconv"
)

var (
	ErrorNotLogin = NewError(1, "请先登录")
)

func NewError(code int, text string) *CodeError {
	return &CodeError{code, text, nil}
}

func NewErrorMsg(text string) *CodeError {
	return &CodeError{0, text, nil}
}

func NewErrorData(code int, text string, data interface{}) *CodeError {
	return &CodeError{code, text, data}
}

func FromError(err error) *CodeError {
	if err == nil {
		return nil
	}
	return &CodeError{0, err.Error(), nil}
}

type CodeError struct {
	Code    int
	Message string
	Data    interface{}
}

func (e *CodeError) Error() string {
	return strconv.Itoa(e.Code) + ": " + e.Message
}
