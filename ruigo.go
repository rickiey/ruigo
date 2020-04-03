package ruigo

import (
	"fmt"
	"log"
	"net/http"
	"regexp"

	json "github.com/json-iterator/go"
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

func (c *Context) JSON(statusCode int, v interface{}) {
	c.Response.WriteHeader(statusCode)
	b, err := json.Marshal(v)
	if err != nil {
		c.String(500, "")
	}
	c.Response.Write(b)
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
	log.Println("[INFO]: Rrequest path :", key)
	for i, handlefunc := range r.handle {
		t, err := regexp.MatchString(i, key)
		if err != nil || !t {
			continue
		}
		handlefunc(c)
		return
	}

	c.String(http.StatusNotFound, fmt.Sprintf("404 NOT FOUND: %s\n", c.Path))

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
