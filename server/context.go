package server

import (
	"encoding/json"
	"io"
	"net/http"
)

type Context struct {
	W          http.ResponseWriter
	R          *http.Request
	PathParams map[string]string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		W:          w,
		R:          r,
		PathParams: map[string]string{},
	}
}
func (c *Context) ReadJson(data interface{}) error {
	body, err := io.ReadAll(c.R.Body)
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
	_, err = c.W.Write(bs)
	if err != nil {
		return err
	}
	c.W.WriteHeader(status)
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

type HandleFunc func(c *Context)

type matchInfo struct {
	n          *node
	pathParams map[string]string
}

func (m *matchInfo) addValue(key string, value string) {
	if m.pathParams == nil {
		// 大多数情况，参数路径只会有一段
		m.pathParams = map[string]string{key: value}
	}
	m.pathParams[key] = value
}
