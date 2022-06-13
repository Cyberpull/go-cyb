package objects

import (
	"reflect"

	"cyberpull.com/go-cyb/errors"
)

type Getter interface {
	Get(key string) (value any, exists bool)
}

func Get[T any](g Getter, key string) (value T, err error) {
	var ok bool
	var tmpValue any

	if tmpValue, ok = g.Get(key); !ok {
		err = errors.Newf(`Key "%s" does not exist`, 500, key)
		return
	}

	if value, ok = tmpValue.(T); !ok {
		rType := reflect.TypeOf(value)
		err = errors.Newf(`Value is not of type "%s"`, rType.Elem().String())
	}

	return
}
