package server

import (
	"github.com/springstar/robot/core"
)

type MatchType int

const (
	MATCH_TYPE_NONE MatchType = iota
	MATCH_TYPE_3V3 
	MATCH_TYPE_5V5

	MATCH_TYPE_10v10_OCCUPY = 4
	MATCH_TYPE_10v10_RESOURCES = 5
	MATCH_TYPE_10V10_CAPTURE = 6
	MATCH_TYPE_REP = 7

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

	if m.isPVPMatch(MatchType(t)) {

	}

	return EXEC_COMPLETED
}

func (m *MatchExecutor) isPVPMatch(t MatchType) bool {
	if t == MATCH_TYPE_3V3 || t == MATCH_TYPE_5V5 {
		return true
	}

	return false
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


