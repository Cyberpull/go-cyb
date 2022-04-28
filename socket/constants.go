package socket

import "time"

const (
	ClientRetry time.Duration = 20

	ErrorPrefix string = "ERROR_MSG::"
	ErrorRcpt   string = "ERROR_RCPT"

	ResponsePrefix string = "RESPONSE::"
	UpdatePrefix   string = "UPDATE::"
)
