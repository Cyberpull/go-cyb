package hooks

import (
	"reflect"
)

var actionHooks = make(map[string][]reflect.Value)

func AddAction(channel string, callback Callback) (err error) {
	fn, err := toActionCallback(callback)

	if err != nil {
		return
	}

	if _, ok := actionHooks[channel]; !ok {
		actionHooks[channel] = make([]reflect.Value, 0)
	}

	actionHooks[channel] = append(actionHooks[channel], fn)

	return
}

func DoActions(channel string, args ...any) (err error) {
	if actions, ok := actionHooks[channel]; ok {
		for _, fn := range actions {
			var value []reflect.Value

			if value, err = call(fn, args...); err != nil {
				break
			}

			if err, _ = value[0].Interface().(error); err != nil {
				break
			}
		}
	}

	return
}

func HasAction(channel string) bool {
	_, ok := actionHooks[channel]
	return ok
}
