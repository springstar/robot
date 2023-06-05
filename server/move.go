package server

import (
	"github.com/springstar/robot/pb"
	"github.com/springstar/robot/msg"
	"fmt"
	_ "math"
	"github.com/springstar/robot/core"
)

type RobotMovement struct {
	r *Robot
	lastSyncTime int64
}

func newMovement(r *Robot) *RobotMovement {
	return &RobotMovement{
		r: r,
	}
}

func (m *RobotMovement) sendMoveRequest(path []*core.Vec2) {
	if len(path) == 0 {
		return
	}

	now := core.GetCurrentTime()
	if !m.isTimeToSync(now) {
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

	msg := msg.SerializeCSStageMove(msg.MSG_CSStageMove, m.r.humanId, start, end, dir, now)
	m.r.sendPacket(msg)

	m.updateSyncTime(now)

}

func (m *RobotMovement) updateSyncTime(now int64) {
	m.lastSyncTime = now
}

func (m *RobotMovement) isTimeToSync(now int64) bool {
	if m.lastSyncTime == 0 || m.lastSyncTime + 100 > now {
		return true
	}

	return false
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
	v := core.MoveTowards(m.r.pos, d, delta)
	if v.Equals(d) {
		return 0
	} else {
		fmt.Println(v)
	}

	m.r.pos = v

	return -1
}


