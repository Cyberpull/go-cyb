package validator

import (
	"regexp"

	"cyberpull.com/go-cyb/errors"
)

var (
	emailRegex1 = regexp.MustCompile("^[a-zA-Z0-9_.]+@[a-zA-Z0-9.]+.[a-zA-Z0-9]$")
	emailRegex2 = regexp.MustCompile("(^[^a-zA-Z0-9]+|[^a-zA-Z0-9]{2,}|[^a-zA-Z0-9]+$)")
)

func checkEmail(data string, info *Validation) (err error) {
	var ok bool

	defer func() {
		if !ok {
			message := message(Email, info)
			err = errors.New(message, 400)
		}
	}()

	if ok = emailRegex1.MatchString(data); !ok {
		return
	}

	ok = !emailRegex2.MatchString(data)

	return
}
