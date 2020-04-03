package ruigo

import (
	"fmt"
	"net/http"
	"regexp"
)

type Context struct {
	Response   http.ResponseWriter
	Request    *http.Request
	KValue     map[interface{}]interface{}
	Path       string
	StatusCode int
	Method     string
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	kv := make(map[interface{}]interface{})
	return &Context{
		Request:  r,
		Response: w,
		KValue:   kv,
		Path:     r.URL.Path,
		Method:   r.Method,
	}
}

func (c *Context) String(statusCode int, msg string) {
	c.StatusCode = statusCode
	c.Response.WriteHeader(statusCode)
	c.Response.Write([]byte(msg))
}

type HandleFunc func(*Context)

type route struct {
	handle map[string]HandleFunc
}

func newRoute() *route {
	m := make(map[string]HandleFunc)
	return &route{m}
}

func (r *route) Add(method, path string, handle HandleFunc) {
	k := method + "-" + path
	if _, ok := r.handle[k]; ok {
		panic("route :  path conflict :" + path)
	}
	r.handle[k] = handle
}

func (r *route) handler(c *Context) {
	key := c.Method + "-" + c.Path
	for i, handlefunc := range r.handle {
		t, err := regexp.Match(i, []byte(key))
		if err != nil || !t {
			c.String(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", c.Path))
		}
		handlefunc(c)
	}
	if handler, ok := r.handle[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", c.Path))
	}
}

type Server struct {
	*route
}

func NewServer() *Server {
	return &Server{newRoute()}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	s.route.handler(c)
}
