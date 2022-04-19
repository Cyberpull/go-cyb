package otp

import (
	_ "cyberpull.com/go-cyb/env"

	"github.com/pquerna/otp"
)

type Digits otp.Digits

const (
	DigitsSix   Digits = Digits(otp.DigitsSix)
	DigitsEight Digits = Digits(otp.DigitsEight)
)

type Options struct {
	Issuer  string
	Account string
	Secret  string
}
