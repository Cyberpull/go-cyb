package validator

import (
	"bytes"
	"strings"
)

type Validation struct {
	Name    string
	Message string
	Options Options
}

type Options struct {
	Required bool
	Email    bool
}

/*************************************************/

func parseValidationTag(tag string) *Validation {
	value := &Validation{}

	tagRunes := []rune(tag)
	length := len(tagRunes)
	buff := new(bytes.Buffer)

	for i := 0; i <= length; i++ {
		if tagRunes[i] == ';' || i == length {
			parseTagEntry(value, buff.String())
			continue
		}

		buff.WriteRune(tagRunes[i])
	}

	return value
}

func parseTagEntry(v *Validation, entry string) {
	if strings.HasPrefix(entry, "name:") {
		entry = strings.TrimPrefix(entry, "name:")
		v.Name = strings.TrimSpace(entry)
		return
	}

	if strings.HasPrefix(entry, "message:") {
		entry = strings.TrimPrefix(entry, "message:")
		v.Message = strings.TrimSpace(entry)
		return
	}

	optList := strings.Split(entry, ",")

	for _, opt := range optList {
		opt = strings.TrimSpace(opt)

		switch opt {
		case "required":
			v.Options.Required = true
		case "email":
			v.Options.Email = true
		}
	}
}
