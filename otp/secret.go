package otp

import (
	"cyberpull.com/go-cyb/errors"

	"github.com/pquerna/otp/totp"
)

func GenerateSecretKey(opts *Options) (err error) {
	if opts.Issuer == "" {
		err = errors.New(`"Issuer" is required`)
		return
	}

	if opts.Account == "" {
		err = errors.New(`"Account" is required`)
		return
	}

	options := totp.GenerateOpts{
		Issuer:      opts.Issuer,
		AccountName: opts.Account,
	}

	key, err := totp.Generate(options)

	if err != nil {
		return
	}

	opts.Secret = key.Secret()

	return
}
