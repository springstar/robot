package core

import (
	"math/rand"
	"time"
)

// gen num betwween [0, n)
func GenRandomInt(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func GenRandomIntList(n int) []int {
	rand.Seed(time.Now().UnixNano())
	a := make([]int, n)
	for i := 0; i < n-1; i++ {
		a[i] = rand.Intn(n)
	}

	return a
}