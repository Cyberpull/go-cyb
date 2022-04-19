package otp

import (
	"fmt"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTP struct {
	opts *Options
}

func (t *TOTP) Generate(ttime time.Time, digits Digits, seconds ...uint) (string, error) {
	if len(seconds) == 0 {
		seconds = append(seconds, 30)
	}

	return totp.GenerateCodeCustom(t.opts.Secret, ttime, totp.ValidateOpts{
		Period:    seconds[0],
		Skew:      1,
		Digits:    otp.Digits(digits),
		Algorithm: otp.AlgorithmSHA1,
	})
}

func (t *TOTP) Validate(code string, digits Digits, seconds ...uint) (bool, error) {
	if len(seconds) == 0 {
		seconds = append(seconds, 30)
	}

	return totp.ValidateCustom(
		code,
		t.opts.Secret,
		time.Now().UTC(),
		totp.ValidateOpts{
			Period:    seconds[0],
			Skew:      1,
			Digits:    otp.Digits(digits),
			Algorithm: otp.AlgorithmSHA1,
		},
	)
}

func (t *TOTP) ToURL() string {
	return fmt.Sprintf(
		"otpauth://%s/%s:%s?secret=%s&issuer=%s",
		"totp",
		t.opts.Issuer,
		t.opts.Account,
		t.opts.Secret,
		t.opts.Issuer,
	)
}

func (t *TOTP) QRCode() *QRCode {
	return NewQR(t.ToURL())
}

/*********************************************/

func NewTOTP(opts Options) *TOTP {
	return &TOTP{
		opts: &opts,
	}
}
