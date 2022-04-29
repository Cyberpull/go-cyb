package validator

import (
	"reflect"

	"cyberpull.com/go-cyb/errors"
)

func Validate(data any, info ...*Validation) (err error) {
	rValue := reflect.ValueOf(data)
	rType := rValue.Type()

	if rType.Kind() == reflect.Pointer {
		rValue = rValue.Elem()
	}

	rKind := rType.Kind()

	if rKind == reflect.Struct {
		err = validateStruct(rValue)
		return
	}

	if len(info) == 0 {
		info = append(info, &Validation{})
	}

	switch rKind {
	case reflect.String:
		err = validateString(rValue.String(), info[0])
	case reflect.Bool:
		err = validateBool(rValue.Bool(), info[0])
	}

	return
}

func validateStruct(data reflect.Value) (err error) {
	length := data.NumField()

	for i := 0; i < length; i++ {
		field := data.Field(i)
		rType := data.Type().Field(i)

		tag, ok := rType.Tag.Lookup("validator")

		if !ok || tag == "" {
			continue
		}

		info := parseValidationTag(rType, tag)

		if err = Validate(field.Interface(), info); err != nil {
			break
		}
	}

	return
}

func validateString(data string, info *Validation) (err error) {
	if info.Options.Required.Value && data == "" {
		message := message(Required, info)
		err = errors.New(message, 400)
		return
	}

	if info.Options.Email.Value {
		err = checkEmail(data, info)

		if err != nil {
			return
		}
	}

	return
}

func validateBool(data bool, info *Validation) (err error) {
	if info.Options.Required.Value && !data {
		message := message(Required, info)
		err = errors.New(message, 400)
	}

	return
}
