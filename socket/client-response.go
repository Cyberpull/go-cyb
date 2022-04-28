package socket

import (
	"context"
	"time"

	"cyberpull.com/go-cyb/errors"
)

type ClientResponseCollection struct {
	mapper map[string]*Output
}

func (c *ClientResponseCollection) Set(uuid string, out *Output) (err error) {
	if _, ok := c.mapper[uuid]; ok {
		err = errors.New("UUID already exists")
		return
	}

	c.mapper[uuid] = out

	return
}

func (c *ClientResponseCollection) Get(req *Request, timeout ...time.Duration) (out *Output, err error) {
	if len(timeout) == 0 {
		timeout = append(timeout, time.Second*10)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), timeout[0])

	defer cancel()

	for {
		select {
		case <-ctx.Done():
			err = errors.New("Request timed out", 408)
			return

		default:
			if tmpOut, ok := c.mapper[req.UUID]; ok {
				delete(c.mapper, req.UUID)
				out = tmpOut
				return
			}

			continue
		}
	}
}

/****************************************/

func newClientResponseCollection() *ClientResponseCollection {
	return &ClientResponseCollection{
		mapper: make(map[string]*Output),
	}
}
