package sys_test

import (
	"fmt"
	"testing"

	"github.com/oswaldoooo/cmicro/sys"
)

func counthello(i int16, name string) {
	fmt.Println(i, name)
}
func jimhello(i int16, name string) {
	fmt.Println(i*2+1, name)
}
func TestTernaryExp(t *testing.T) {
	sys.TernaryExpressionFunc(1 > 2, counthello, jimhello, 23, "jesko")
	sys.TernaryExpressionFunc(1 < 2, counthello, jimhello, 23, "jesko")
}

func TestTernaryExpFunc(t *testing.T) {
	sys.TernaryExpressFunc(1 < 2, sys.ToFunc(counthello, 23, "jesko"), sys.ToFunc(jimhello, 23, "jim"))
}
