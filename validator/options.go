package validator

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
