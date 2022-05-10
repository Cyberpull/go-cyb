package socket

func write[T any](out chan T, data T) {
	defer recover()

	if out == nil {
		return
	}

	out <- data
}
