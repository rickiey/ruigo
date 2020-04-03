package ruigo

import (
	"fmt"
	"regexp"
	"testing"
)

func Test_re(t *testing.T) {
	s := "GET-/he*"
	fmt.Println(regexp.MatchString(s, "GET-/hello"))
}
