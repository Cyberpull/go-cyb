package errors

import (
	"errors"
	"fmt"
	"runtime/debug"
)

type Error struct {
	code    int
	message string
	stack   string
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Error() string {
	message := e.message

	if e.stack != "" {
		message += "\n" + e.stack
	}

	return message
}

func (e *Error) WithStack() *Error {
	e.stack = string(debug.Stack())
	return e
}

/**********************************/

func New(message string, code ...int) *Error {
	code = sanitizeErrorCode(code...)

	return &Error{
		code:    code[0],
		message: message,
	}
}

func Newf(message string, v ...any) *Error {
	code := make([]int, 0)

	if len(v) > 0 {
		vcode, ok := v[0].(int)

		if ok {
			code = append(code, vcode)
		}

		v = v[1:]
	}

	message = fmt.Sprintf(message, v...)

	return New(message, code...)
}

func From(v any, code ...int) *Error {
	var value *Error

	switch x := v.(type) {
	case *Error:
		value = x

		if len(code) > 0 {
			code = sanitizeErrorCode(code...)
			value.code = code[0]
		}
	case string:
		value = New(x, code...)
	case error:
		value = New(x.Error(), code...)
	default:
		value = New("An unknown error occurred", 500)
	}

	return value
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}

func sanitizeErrorCode(code ...int) []int {
	if len(code) == 0 {
		code = append(code, 500)
	}

	if code[0] == 0 || code[0] == 200 {
		code[0] = 500
	}

	return code
}
