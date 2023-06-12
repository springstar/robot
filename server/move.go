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
	path []*core.Vec2
	lastSyncTime int64
}

func newMovement(r *Robot) *RobotMovement {
	return &RobotMovement{
		r: r,
	}
}

func (m *RobotMovement) sendMoveRequest() {
	if len(m.path) == 0 {
		return
	}

	now := core.GetCurrentTime()

	start := &pb.DVector3{
		X: m.r.pos.X,
		Y: m.r.pos.Y,
		Z: 0,
	}

	var end []*pb.DVector3
	for _, p := range m.path {
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
	
	for _, pos := range m.path {
		fmt.Println(pos)
	}
	// clear path
	m.path = m.path[:0]
}

func (m *RobotMovement) updateSyncTime(now int64) {
	m.lastSyncTime = now
}

func (m *RobotMovement) isTimeToSync(now int64) bool {
	if m.lastSyncTime == 0 || m.lastSyncTime + 100 < now {
		return true
	}

	return false
}

func (m *RobotMovement) exec(params []string, delta int) ExecState {
	v, err := core.Str2Int(params[0])
	if err != nil {
		fmt.Printf("path point error %s", params)
		return EXEC_COMPLETED
	}

	target := serv.sceneMgr.getPoint(m.r.mapSn, v)

	// fmt.Printf("map %d num %d target %v\n", m.r.mapSn, v, target)
	delta = 5
	// delta = delta * int(m.r.speed)
	now := core.GetCurrentTime()

	r := m.moveto(target, float32(delta))

	if m.isTimeToSync(now) {
		m.sendMoveRequest()
	}

	if r == -1 {
		return EXEC_ONGOING
	}

	return EXEC_COMPLETED
}

func (m *RobotMovement)moveto(d *core.Vec2, delta float32) int {
	// fmt.Println(m.r.pos, *d, delta)
	v := core.MoveTowards(m.r.pos, d, delta)
	if v.Equals(d) {
		fmt.Println("move completed ", d)
		return 0
	} else {
		// fmt.Println(v)
	}

	m.r.pos = v
	m.path = append(m.path, v)

	return -1
}


