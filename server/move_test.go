package server

import (
	"github.com/springstar/robot/core"
	_ "fmt"
	"testing"
)

func TestMoveTo(t *testing.T) {
	a := core.NewVec2(212.56, 75.3)
	b := core.NewVec2(100.25, 200.30)


	for i := 0; i < 10000000; i++ {
		v := core.MoveTowards(a, b, 15)
		if v.Equals(b) {
			t.Log("reach destination")
		}
	}

}