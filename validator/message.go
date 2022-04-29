package validator

import "strings"

type OptionName string

const (
	Required OptionName = "Required"
	Email    OptionName = "Email"
)

var messageMapper map[OptionName]string

func init() {
	messageMapper = make(map[OptionName]string)

	messageMapper[Required] = "$name is required"
	messageMapper[Email] = "$name does not appear to be valid"
}

func message(name OptionName, info *Validation) (value string) {
	value = info.Options.message(name)

	if value == "" {
		value = messageMapper[name]
	}

	value = strings.ReplaceAll(value, "$name", info.Name)

	return
}
