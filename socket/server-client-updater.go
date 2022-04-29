package socket

import (
	"cyberpull.com/go-cyb/errors"
	"cyberpull.com/go-cyb/objects"
)

type ServerClientUpdater struct {
	ref *ServerClientRef
}

func (s *ServerClientUpdater) Update(args ...any) (err error) {
	var ok bool

	length := len(args)

	switch length {
	case 1:
		var out *Output

		if out, ok = args[0].(*Output); ok {
			err = s.writeUpdate(out)

			return
		}
	case 3, 4:
		code := make([]int, 0)
		var method, channel string

		if method, ok = args[0].(string); !ok {
			break
		}

		if channel, ok = args[1].(string); !ok {
			break
		}

		if length == 4 {
			var vcode int

			if vcode, ok = args[3].(int); !ok {
				break
			}

			code = append(code, vcode)
		}

		err = s.updateRaw(method, channel, args[2], code...)

		return
	}

	err = errors.New("Invalid arguments")

	return
}

func (s *ServerClientUpdater) updateRaw(method, channel string, data any, code ...int) (err error) {
	if len(code) == 0 {
		code = append(code, 200)
	}

	out := &Output{
		Method:  method,
		Channel: channel,
		Code:    code[0],
	}

	if err = out.SetData(data); err != nil {
		return
	}

	err = s.writeUpdate(out)

	return
}

func (s *ServerClientUpdater) writeUpdate(output *Output) (err error) {
	data := []byte(UpdatePrefix)

	json, err := objects.ToJSON(output)

	if err != nil {
		return
	}

	data = append(data, json...)

	_, err = s.ref.Writeln(data)

	return
}

/*******************************************/

func newServerClientUpdater(ref *ServerClientRef) *ServerClientUpdater {
	return &ServerClientUpdater{ref: ref}
}
