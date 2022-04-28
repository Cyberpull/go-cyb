package socket

import "time"

const (
	ClientRetry time.Duration = 20

	ErrorPrefix string = "ERROR_MSG::"
	ErrorRcpt   string = "ERROR_RCPT"

	ResponseTxt    string = "RESPONSE"
	ResponsePrefix string = ResponseTxt + "::"

	UpdateTxt    string = "UPDATE"
	UpdatePrefix string = UpdateTxt + "::"
)
