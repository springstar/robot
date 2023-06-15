package server

type TeamExecutor struct {
	*Robot
	*TeamMember
}

type TeamMember struct {
	teamId int64
	isLeader bool
}

func newTeamMember(tid int64) *TeamMember {
	return &TeamMember{
		teamId: tid,
	}
}

func newTeamExecutor(r *Robot) *TeamExecutor {
	return &TeamExecutor{
		Robot: r,
	}
}



func (t *TeamExecutor) exec(params []string, delta int) ExecState {
	return EXEC_COMPLETED
}

func (t *TeamExecutor) checkIfExec(params []string) bool {

	return true
}

func (t *TeamExecutor) handleBreak() {
	// now := core.GetCurrentTime()

}