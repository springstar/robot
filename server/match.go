package server

type MatchExecutor struct {
	*Robot
}

func newMatchExecutor(r *Robot) *MatchExecutor {
	return &MatchExecutor{
		Robot: r,
	}
}

func (m *MatchExecutor) exec(params []string, delta int) ExecState {
	return EXEC_COMPLETED
}

func (m *MatchExecutor) checkIfExec() bool {
	return true
}