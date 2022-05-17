package socket

func writeOne[T any](out []chan T, data T) {
	if len(out) > 0 {
		write(out[0], data)
	}
}

func write[T any](out chan T, data T) {
	defer recover()

	if out == nil {
		return
	}

	out <- data
}
