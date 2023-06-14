package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

type MatchType int

const (
	MT_NONE MatchType = iota
	MT_3V3 
	MT_5V5

	MT_OCCUPY = 4
	MT_RESOURCES = 5
	MT_CAPTURE = 6
	MT_REP = 7

)

type MatchExecutor struct {
	*Robot
	isMatching bool
	mt MatchType
	startTime int64
}

func newMatchExecutor(r *Robot) *MatchExecutor {
	return &MatchExecutor{
		Robot: r,
		isMatching: false,
		startTime: 0,
		mt: MT_NONE,
	}
}

func (m *MatchExecutor) exec(params []string, delta int) ExecState {
	t, err := core.Str2Int(params[0])
	if err != nil {
		core.Error("match type error %s", params)
		return EXEC_COMPLETED
	}

	if m.isPVEMatch(MatchType(t)) {
		m.startPVEMatch()
	} else if m.isPVPMatch(MatchType(t)) {
		m.startPVPMatch(MatchType(t))
	} else {
		core.Error("error match type ", t)
	}

	return EXEC_COMPLETED
}

func (m *MatchExecutor) isPVEMatch(t MatchType) bool {
	if t == MT_REP {
		return true
	}

	return false
}

func (m *MatchExecutor) isPVPMatch(t MatchType) bool {
	if t == MT_3V3 || t == MT_5V5 || t == MT_OCCUPY || t == MT_RESOURCES || t == MT_CAPTURE {
		return true
	} else {
		return false
	}
}

func (m *MatchExecutor) startPVPMatch(mt MatchType) {
	core.Info("start pvp match")
	request := msg.SerializeCSMatchEnrollRequest(msg.MSG_CSMatchEnrollRequest, int32(mt))
	m.sendPacket(request)
	m.markMatching(mt)
}

func (m *MatchExecutor) startPVEMatch() {

}

func (m *MatchExecutor) markMatching(t MatchType) {
	m.mt = t
	m.isMatching = true
	m.startTime = core.GetCurrentTime()
}

func (m *MatchExecutor) markNone() {
	m.mt = MT_NONE
	m.isMatching = false
	m.startTime = 0
}

func (m *MatchExecutor) cancelArenaMatch() {
	request := msg.SerializeCSMatchCancelRequest(msg.MSG_CSMatchCancelRequest, int32(m.mt))
	m.sendPacket(request)
}

func (m *MatchExecutor) checkIfExec() bool {
	// now := core.GetCurrentTime()
	if m.isMatching {
		return false
	}

	return true
}

func (m *MatchExecutor) handleBreak() {
	// now := core.GetCurrentTime()

}

func (r *Robot) handleArenaEnroll(packet *core.Packet) {
	msg := msg.ParseSCMatchEnrollResponse(int32(msg.MSG_SCMatchEnrollResponse), packet.Data)
	code := msg.GetCode()
	if code == 1 {
		executor := r.findExecutor("match").(*MatchExecutor)
		executor.markNone()
	}

	core.Info("arena enroll code ", code)
}

func (r *Robot) handleArenaMatchResult(packet *core.Packet) {
	dumpMatchResult(packet)
}

func dumpMatchResult(packet *core.Packet) {
	msg := msg.ParseSCMatchResult(int32(msg.MSG_SCMatchResult), packet.Data)
	sideA := msg.GetSideA()
	sideB := msg.GetSideB()

	for _, member := range sideA.GetMembers() {
		core.Info(member.GetName())
	}

	for _, member := range sideB.GetMembers() {
		core.Info(member.GetName())
	}
}
