package chomolungma

import (
	"encoding/json"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	// origin objects
	Res http.ResponseWriter
	Req *http.Request
	// request data
	Path   string
	Method string
	Params map[string]string
	// response data
	StatusCode int
	// middleware
	handlers []HandlerFunc
	index    int
}

func newContext(res http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Res:    res,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:  -1,
	}
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Res.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Res.Header().Set(key, value)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Res)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Res, err.Error(), 500)
	}
}
