package core

import (
	"math/rand"
	"time"
	_ "math"
)

type Vector2 struct {
	X float32
	Y float32
}

type Vector3 struct {
	x float32
	y float32
	z float32
}

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

// func Rotation(p1x, p1y, p2x, p2y float64) float64 {
// 	return (float64)(math.Atan2(p1x - p2x, p2y - p1y) * 180.0 / math.PI + 630) % 360.0;
// }