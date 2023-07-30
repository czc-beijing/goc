package server

import (
	"goc/context"
)

func Hello(c *context.Context) {
	d := map[string]string{
		"h": "hello",
	}
	c.OkJson(d)
	return
}
