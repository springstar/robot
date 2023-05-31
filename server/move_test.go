package server

import (
	"fmt"
	"testing"
)

func TestMoveTo(t *testing.T) {
	a := newObj()
	
	a.x = 12.56
	a.y = 7.3

	b := newObj()
	b.x = 100.25
	b.y = 200.30

	for i := 0; i < 10000000; i++ {
		moveto(a, b)
	}

	fmt.Println(a.x - b.x)
	fmt.Println(a.y - b.y)
}