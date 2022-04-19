package errors

import (
	"errors"
	"fmt"
)

type Error struct {
	code    int
	message string
}

func (e Error) Code() int {
	return e.code
}

func (e Error) Error() string {
	return e.message
}

/**********************************/

func New(message string, code ...int) *Error {
	if len(code) == 0 {
		code = append(code, 500)
	}

	if code[0] == 0 {
		code[0] = 500
	}

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
