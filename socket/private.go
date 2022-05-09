package socket

func write[T any](out chan T, data T) {
	if out != nil {
		out <- data
	}
}
