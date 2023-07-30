package context

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		w: w,
		r: r,
	}
}
func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.r.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, data)
}

func (c *Context) WriteJson(status int, data interface{}) error {
	bs, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.w.Write(bs)
	if err != nil {
		return err
	}
	c.w.WriteHeader(status)
	return err
}

func (c *Context) OkJson(data interface{}) error {
	return c.WriteJson(http.StatusOK, data)
}

func (c *Context) ErrorJson(data interface{}) error {
	return c.WriteJson(http.StatusInternalServerError, data)
}

func (c *Context) BadRequestJson(data interface{}) error {
	return c.WriteJson(http.StatusBadRequest, data)
}
