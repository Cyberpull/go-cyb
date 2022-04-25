package hooks

import (
	"reflect"

	"cyberpull.com/go-cyb/errors"
)

type Callback any

func toCallback(fn Callback) (value reflect.Value, err error) {
	tmpValue := reflect.ValueOf(fn)

	if tmpValue.Kind() != reflect.Func {
		err = errors.Newf(`"%s" is not a function`, 500, tmpValue.Type().Name())
		return
	}

	value = tmpValue

	return
}

func toActionCallback(fn Callback) (value reflect.Value, err error) {
	tmpValue, err := toCallback(fn)

	if err != nil {
		return
	}

	fnType := tmpValue.Type()

	if fnType.NumOut() != 1 {
		err = errors.New("Callback function must have only 1 return value")
		return
	}

	errType := reflect.TypeOf((*error)(nil)).Elem()

	if !fnType.Out(0).Implements(errType) {
		err = errors.New(`Callback's return value must be an "error".`, 500)
		return
	}

	value = tmpValue

	return
}

func toFilterCallback(fn Callback) (value reflect.Value, err error) {
	tmpValue, err := toCallback(fn)

	if err != nil {
		return
	}

	fnType := tmpValue.Type()

	if fnType.NumOut() != 2 {
		err = errors.New("Callback function must have 2 return values")
		return
	}

	errType := reflect.TypeOf((*error)(nil)).Elem()

	if !fnType.Out(1).Implements(errType) {
		err = errors.New(`Callback's last return value must be an "error".`, 500)
		return
	}

	value = tmpValue

	return
}

func call(fn reflect.Value, args ...any) (value []reflect.Value, err error) {
	fargs := make([]reflect.Value, len(args))

	for i, arg := range args {
		fargs[i] = reflect.ValueOf(arg)
	}

	defer func() {
		if r := recover(); r != nil {
			err = errors.From(r)
		}
	}()

	value = fn.Call(fargs)

	return
}
