package validator

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strings"
)

func parseValidationTag(field reflect.StructField, tag string) *Validation {
	value := &Validation{}

	tagRunes := []rune(tag)
	length := len(tagRunes)
	buff := new(bytes.Buffer)

	for i := 0; i <= length; i++ {
		if i == length || tagRunes[i] == ';' {
			parseValidationTagEntry(value, buff.String())
			buff.Reset()
			continue
		}

		buff.WriteRune(tagRunes[i])
	}

	sanitizeValidation(field, value)

	return value
}

func parseValidationTagEntry(v *Validation, entry string) {
	chunks := strings.SplitN(entry, ":", 2)

	switch chunks[0] {
	case "required":
		parseValidationTagEntryOptionValue(&v.Options.Required, chunks)

	case "email":
		parseValidationTagEntryOptionValue(&v.Options.Email, chunks)

	case "fieldName":
		if len(chunks) == 2 {
			v.Name = chunks[1]
		}
	}
}

func parseValidationTagEntryOptionValue(v *OptionValue, chunks []string) {
	var err error

	switch len(chunks) {
	case 1:
		v.Value = true

	case 2:
		chunks2 := strings.Split(chunks[1], "|")

		for _, valueString := range chunks2 {
			var value bool

			err = json.Unmarshal([]byte(valueString), &value)

			if err == nil {
				v.Value = value
				continue
			}

			prefix := "message="

			if strings.HasPrefix(valueString, prefix) {
				v.Message = strings.TrimPrefix(valueString, prefix)
			}
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
