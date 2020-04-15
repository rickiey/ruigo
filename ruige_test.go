package ruigo

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"testing"
)

func Test_re(t *testing.T) {

	path := "/asd545"
	fmt.Println(regexp.MatchString("^/[A-Za-z1-9/_-]+$", path))

	start()
}

func start() {
	s := NewServer()
	s.Add([]string{"POST"}, "/hello", Hello)
	port := ":9000"
	log.Println("[INFO] Start http server :", port)
	log.Fatal(http.ListenAndServe(port, s))
}

func Hello(c *Context) {
	c.JSON(200, map[string]int{
		"code": 200,
	})
}
