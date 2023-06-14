package server

import (
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
	startTime int64
}

func newMatchExecutor(r *Robot) *MatchExecutor {
	return &MatchExecutor{
		Robot: r,
		isMatching: false,
		startTime: 0,
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
		m.startPVPMatch()
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

func (m *MatchExecutor) startPVPMatch() {

}

func (m *MatchExecutor) startPVEMatch() {

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


