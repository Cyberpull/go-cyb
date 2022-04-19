package hooks

import (
	"reflect"

	"cyberpull.com/go-cyb/errors"
)

var filterHooks = make(map[string][]reflect.Value)

func AddFilter(channel string, callback Callback) {
	fn, err := toFilterCallback(callback)

	if err != nil {
		return
	}

	if _, ok := filterHooks[channel]; !ok {
		filterHooks[channel] = make([]reflect.Value, 0)
	}

	filterHooks[channel] = append(filterHooks[channel], fn)
}

func ApplyFilters[T any](channel string, data T, args ...interface{}) (value T, err error) {
	value = data

	if filters, ok := filterHooks[channel]; ok {
		for _, fn := range filters {
			args[0] = value

			var retValues []reflect.Value

			if retValues, err = call(fn, args...); err != nil {
				break
			}

			if err, _ = retValues[1].Interface().(error); err != nil {
				break
			}

			var ok bool

			if value, ok = retValues[0].Interface().(T); !ok {
				err = errors.New(`Invalid return value`)
				break
			}
		}
	}

	return
}

func HasFilter(channel string) bool {
	_, ok := filterHooks[channel]
	return ok
}
