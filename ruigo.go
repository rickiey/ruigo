package ruigo

import (
	json "github.com/json-iterator/go"
	"log"
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

func (c *Context) JSON(statusCode int, v interface{}) {
	c.Response.WriteHeader(statusCode)
	b, err := json.Marshal(v)
	if err != nil {
		c.String(500, "")
	}
	c.Response.Write(b)
}

type HandleFunc func(*Context)

type HandleSet struct {
	method []string
	hd     HandleFunc
}

type route struct {
	handle map[string]HandleSet
}

func newRoute() *route {
	m := make(map[string]HandleSet)
	return &route{m}
}

var METHODS = map[string]bool{
	http.MethodGet:     true,
	http.MethodPost:    true,
	http.MethodPatch:   true,
	http.MethodDelete:  true,
	http.MethodConnect: true,
	http.MethodHead:    true,
	http.MethodOptions: true,
	http.MethodPut:     true,
	http.MethodTrace:   true,
}

func (r *route) Add(method []string, path string, handle HandleFunc) {
	for _, v := range method {
		if _, ok := METHODS[v]; !ok {
			log.Panicf("unsupport http method %v ; it is one of  GET POST PUT DELETE HEAD TRACE OPTIONS CONNECT PATCH \n", v)
		}
	}
	ok, err := regexp.MatchString("^/[A-Za-z1-9/_-]+$", path)
	if err != nil {
		log.Panicln("match path error :", err)
	}
	if !ok {
		log.Panicln("path not in compliance : ", path)
	}
	k := HandleSet{
		method: method,
		hd:     handle,
	}
	if _, ok := r.handle[path]; ok {
		panic("route :  path conflict :" + path)
	}
	r.handle[path] = k
}

func (r *route) handler(c *Context) {

	hs, ok := r.handle[c.Path]
	if !ok {
		c.Response.WriteHeader(http.StatusNotFound)
		return
	}
	for _, v := range hs.method {
		if c.Method == v {
			hs.hd(c)
		}
	}
	c.Response.WriteHeader(http.StatusMethodNotAllowed)
	return

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
