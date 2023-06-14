package server

import (
	_ "github.com/springstar/robot/core"
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
	return EXEC_COMPLETED
}

func (m *MatchExecutor) checkIfExec() bool {
	// now := core.GetCurrentTime()
	if m.isMatching {

		return false
	}


	return true
}

func (m *MatchExecutor) handleBreak() {
	
}


