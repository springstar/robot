package server

type TeamExecutor struct {
	*Robot

}

func newTeamExecutor(r *Robot) *TeamExecutor {
	return &TeamExecutor{
		Robot: r,
	}
}

func (t *TeamExecutor) exec(params []string, delta int) ExecState {
	return EXEC_COMPLETED
}

func (t *TeamExecutor) checkIfExec() bool {
	return true
}

func (t *TeamExecutor) handleBreak() {
	// now := core.GetCurrentTime()

}