package server

import (
	"math"
)

type Obj struct {
	x float64
	y float64
}

func newObj() *Obj {
	return &Obj{}
}

func moveto(a *Obj, b *Obj) {
	dir := rotation(a.x, a.y, b.x, b.y)
	a.x += math.Cos(dir)*0.5
	a.y += math.Sin(dir)*0.5
	
}

func rotation(srcX, srcY float64, dstX, dstY float64) float64 {
	return math.Atan2(dstY - srcY, dstX - srcX)
}