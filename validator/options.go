package validator

import (
	"bytes"
	"reflect"
	"strings"
)

type Validation struct {
	Name    string
	Options Options
}

type Options struct {
	Required OptionValue
	Email    OptionValue
}

func (o Options) message(name OptionName) string {
	switch name {
	case Required:
		return o.Required.Message
	case Email:
		return o.Email.Message
	}

	return ""
}

type OptionValue struct {
	Value   bool
	Message string
}

/*************************************************/

func parseValidationTag(field reflect.StructField, tag string) *Validation {
	value := &Validation{}

	tagRunes := []rune(tag)
	length := len(tagRunes)
	buff := new(bytes.Buffer)

	for i := 0; i <= length; i++ {
		if i == length || tagRunes[i] == ';' {
			parseTagEntry(value, buff.String())
			buff.Reset()
			continue
		}

		buff.WriteRune(tagRunes[i])
	}

	sanitizeValidation(field, value)

	return value
}

func parseTagEntry(v *Validation, entry string) {
	if strings.HasPrefix(entry, "name:") {
		entry = strings.TrimPrefix(entry, "name:")
		v.Name = strings.TrimSpace(entry)
		return
	}

	// if strings.HasPrefix(entry, "message:") {
	// 	entry = strings.TrimPrefix(entry, "message:")
	// 	v.Message = strings.TrimSpace(entry)
	// 	return
	// }

	optList := strings.Split(entry, ",")

	for _, opt := range optList {
		opt = strings.TrimSpace(opt)

		switch opt {
		case "required":
			v.Options.Required.Value = true
		case "email":
			v.Options.Email.Value = true
		}
	}
}

func sanitizeValidation(field reflect.StructField, v *Validation) {
	if v.Name == "" {
		v.Name = field.Name
	}

	sanitizeOptions(field, &v.Options)
}

func sanitizeOptions(field reflect.StructField, opts *Options) {
	// ...
}
