package server

import (
	"fmt"
	_ "math"
	"github.com/springstar/robot/core"
)

type RobotMovement struct {
	*core.Vec2
	r *Robot
}

func newMovement(r *Robot) *RobotMovement {
	return &RobotMovement{
		Vec2: core.NewVec2(0, 0),
		r: r,
	}
}

func (m *RobotMovement) exec(params []string, delta int) ExecState {
	var target []float64
	for _, para := range params {
		v := core.Str2Float64(para)
		target = append(target, v)
	}

	d := core.NewVec2(float32(target[0]), float32(target[1]))
	
	delta = delta * int(m.r.speed)
	r := m.moveto(d, float32(delta))

	if r == -1 {
		return EXEC_ONGOING
	}

	return EXEC_COMPLETED
}

func (m *RobotMovement)moveto(d *core.Vec2, delta float32) int {
	v := core.MoveTowards(m.Vec2, d, delta)
	if v.Equals(d) {
		return 0
	} else {
		fmt.Println(v)
	}

	m.Vec2 = v

	return -1
}


