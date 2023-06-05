package server

import (
	"github.com/springstar/robot/pb"
	"github.com/springstar/robot/msg"
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

func (m *RobotMovement) sendMoveRequest(path []*core.Vec2) {
	if len(path) == 0 {
		return
	}

	start := &pb.DVector3{
		X: m.r.pos.X,
		Y: m.r.pos.Y,
		Z: 0,
	}

	var end []*pb.DVector3
	for _, p := range path {
		pos := &pb.DVector3{
			X: p.X,
			Y: p.Y,
			Z: 0,
		}

		end = append(end, pos)
	}

	dir := &pb.DVector3{
		X: m.r.dir.X,
		Y: m.r.dir.Y,
		Z: 0,
	}

	now := core.GetCurrentTime()

	msg := msg.SerializeCSStageMove(msg.MSG_CSStageMove, m.r.humanId, start, end, dir, now)
	m.r.sendPacket(msg)
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


