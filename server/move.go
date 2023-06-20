package server

import (
	 _ "log"
	"github.com/springstar/robot/pb"
	"github.com/springstar/robot/msg"
	_ "math"
	"github.com/springstar/robot/core"
)

type RobotMovement struct {
	r *Robot
}

func newMovement(r *Robot) *RobotMovement {
	return &RobotMovement{
		r: r,
	}
}

func (r *Robot) sendMoveRequest() {
	if len(r.path) == 0 {
		return
	}

	now := core.GetCurrentTime()

	start := &pb.DVector3{
		X: r.lastPos.X,
		Y: r.lastPos.Y,
		Z: 0,
	}

	var end []*pb.DVector3
	for _, p := range r.path {
		pos := &pb.DVector3{
			X: p.X,
			Y: p.Y,
			Z: 0,
		}

		end = append(end, pos)
	}

	dir := &pb.DVector3{
		X: r.dir.X,
		Y: r.dir.Y,
		Z: 0,
	}

	msg := msg.SerializeCSStageMove(msg.MSG_CSStageMove, r.humanId, start, end, dir, now)
	r.sendPacket(msg)
	

	r.updateSyncTime(now)
	
	// clear path
	r.path = r.path[:0]
	r.lastPos = r.pos
}

func (r *Robot) updateSyncTime(now int64) {
	r.lastSyncTime = now
}

func (r *Robot) isTimeToSync(now int64) bool {
	if r.lastSyncTime == 0 || r.lastSyncTime + 150 < now {
		return true
	}

	return false
}

func (m *RobotMovement) checkIfExec(params []string) bool {
	return true
}

func (m *RobotMovement) handleBreak() {

}

func (m *RobotMovement) exec(params []string, delta int) ExecState {
	v, err := core.Str2Int(params[0])
	if err != nil {
		core.Error("path point error %s", params)
		return EXEC_COMPLETED
	}

	// core.Info("move sn ", v)
	target := serv.sceneMgr.getPoint(m.r.mapSn, v)
	r := m.r.move(target)

	if r == -1 {
		return EXEC_REPEATED
	}

	return EXEC_COMPLETED
}

func (r *Robot) move(target *core.Vec2) int{
	now := core.GetCurrentTime()

	delta := 5
	// delta = delta * int(m.r.speed)

	if r.lastPos == nil {
		r.lastPos = r.pos
	}

	ret := r.moveto(target, float32(delta))

	if r.isTimeToSync(now) {
		r.sendMoveRequest()
	}

	return ret
}

func (r *Robot) moveto(d *core.Vec2, delta float32) int {
	v := core.MoveTowards(r.pos, d, delta)
	if v.Equals(d) {
		return 0
	} else {
		// fmt.Println(v)
	}

	r.pos = v
	r.path = append(r.path, v)

	return -1
}


