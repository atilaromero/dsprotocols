package consensus

import (
	"fmt"
	"testing"
)

func TestTmp(testing *testing.T) {
	fmt.Println("teste")
	a := fmt.Sprintf("%s", "ASDF")
	fmt.Println(a)
	var b string
	var c int
	_, err := fmt.Sscanf(a, "%s %d", &b, &c)
	fmt.Print(b, c, err)
}
