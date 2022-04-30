package otp

type Account interface {
	OTPAccount() string
}

func GetAccount(i Account) string {
	return i.OTPAccount()
}
