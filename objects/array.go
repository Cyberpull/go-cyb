package objects

type Array[T comparable] struct {
	data []T
}

func (a *Array[T]) First() T {
	return a.data[0]
}

func (a *Array[T]) Last() T {
	return a.data[a.Length()-1]
}

func (a *Array[T]) At(i int) T {
	return a.data[i]
}

func (a *Array[T]) Get(i int) T {
	return a.At(i)
}

func (a *Array[T]) Take(i int) T {
	value := a.At(i)
	a.Splice(i, 1)
	return value
}

func (a *Array[T]) Push(v ...T) int {
	a.data = append(a.data, v...)
	return a.Length() - 1
}

func (a *Array[T]) Pop() T {
	var value T

	lastIndex := len(a.data)

	if lastIndex > 0 {
		lastIndex -= 1

		value = a.data[lastIndex]
		a.data = a.data[:lastIndex]
	}

	return value
}

func (a *Array[T]) Slice(start int, stop ...int) *Array[T] {
	var value []T

	if len(stop) == 0 {
		value = a.data[start:]
	} else {
		value = a.data[start:stop[0]]
	}

	return NewArray(value...)
}

func (a *Array[T]) Splice(offset int, length int, v ...T) *Array[T] {
	value := make([]T, 0)

	endOffset := offset + length

	start := a.data[:offset]
	value = a.data[offset:endOffset]
	end := a.data[endOffset:]

	a.data = append(start, v...)
	a.data = append(a.data, end...)

	return NewArray(value...)
}

func (a *Array[T]) Contains(v T) bool {
	return a.IndexOf(v) >= 0
}

func (a *Array[T]) IndexOf(v T) int {
	for i, data := range a.data {
		if data == v {
			return i
		}
	}

	return -1
}

func (a *Array[T]) Length() int {
	return len(a.data)
}

/*************************************/

func NewArray[T comparable](data ...T) *Array[T] {
	value := &Array[T]{
		data: data,
	}

	return value
}
