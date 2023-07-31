package server

func Hello(c *Context) {
	d := map[string]string{
		"h": "hello",
	}
	c.OkJson(d)
	return
}
