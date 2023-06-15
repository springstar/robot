package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

type TeamOperation int32

const (
	TO_NONE = iota
	TO_CREATE
)

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
	op, _ := parseOperation(params)
	switch op {
	case TO_CREATE:
		t.create(params)
	}
	return EXEC_COMPLETED
}

func (t *TeamExecutor) create(params []string) {
	targetSn, _ := core.Str2Int(params[1])
	teamType, _ := core.Str2Int(params[2])
	request := msg.SerializeCSPlatCreateTeam(msg.MSG_CSPlatCreateTeam, int32(targetSn), int32(teamType))
	t.sendPacket(request)
}

func parseOperation(params []string) (TeamOperation, error){
	op, err := core.Str2Int(params[0])
	if err != nil {
		return TO_NONE, err
	}

	return TeamOperation(op), nil
}

func (t *TeamExecutor) checkIfExec(params []string) bool {
	op, err := parseOperation(params)
	if err != nil {
		return false
	}

	switch op {
	case TO_CREATE:
		return t.teamId == 0
	default:
		return false	
	}
	return true
}

func (t *TeamExecutor) handleBreak() {
	// now := core.GetCurrentTime()

}