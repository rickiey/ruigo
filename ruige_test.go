package ruigo

import (
	"fmt"
	"regexp"
	"testing"
	"log"
	"net/http"

	"github.com/rickiey/ruigo"
)

func Test_re(t *testing.T) {
	s := "GET-/he*"
	fmt.Println(regexp.MatchString(s, "GET-/hello"))
}

func main() {
	s := ruigo.NewServer()
	s.Add("GET", "/hello*", Hello)
	port := ":9000"
	log.Println("[INFO] Start http server :", port)
	log.Fatal(http.ListenAndServe(port, s))
}

func Hello(c *ruigo.Context) {
	c.JSON(200, map[string]int{
		"code": 200,
	})
}
