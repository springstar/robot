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
	ae *AsyncExecutor
	*Robot
	*TeamMember
}
type TeamMember struct {
	teamId int64
	isLeader bool
}

func newTeamMember(tid int64, isLeader bool) *TeamMember {
	return &TeamMember{
		teamId: tid,
		isLeader: isLeader,
	}
}

func newTeamExecutor(r *Robot) *TeamExecutor {
	return &TeamExecutor{
		Robot: r,
		ae: newAsyncExecutor(),
		TeamMember: newTeamMember(0, false),
	}
}

func (t *TeamExecutor) exec(params []string, delta int) ExecState {
	op, _ := parseOperation(params)
	switch op {
	case TO_CREATE:
		return t.create(params)
	}
	return EXEC_COMPLETED
}

func (t *TeamExecutor) create(params []string) ExecState {
	core.Info("send create team request")

	targetSn, _ := core.Str2Int(params[1])
	teamType, _ := core.Str2Int(params[2])
	request := msg.SerializeCSPlatCreateTeam(msg.MSG_CSPlatCreateTeam, int32(targetSn), int32(teamType))
	t.sendPacket(request)
	t.ae.setOngoing()
	return t.ae.getState()
}

func parseOperation(params []string) (TeamOperation, error){
	op, err := core.Str2Int(params[0])
	if err != nil {
		return TO_NONE, err
	}

	return TeamOperation(op), nil
}

func (t *TeamExecutor) checkIfExec(params []string) bool {
	if !t.ae.checkIfExec(params) {
		return false
	}

	if t.profession == 0 {
		return false
	}

	op, err := parseOperation(params)
	if err != nil {
		return false
	}

	switch op {
	case TO_CREATE:
		return t.teamId <= 0
	default:
		return false	
	}
	return true
}

func (t *TeamExecutor) handleBreak() {
	// now := core.GetCurrentTime()

}

func (r *Robot) handleTeamDetail(packet *core.Packet) {
	detail := msg.ParseSCTeamMine(int32(msg.MSG_SCTeamMine), packet.Data)
	e := r.findExecutor("team").(*TeamExecutor)
	teamId := detail.GetTeamId() 
	if teamId <= 0 {
		return
	}

	isLeader := r.isLeader(detail.GetLeaderId())
	e.TeamMember.teamId = detail.GetTeamId()
	e.TeamMember.isLeader = isLeader
	e.ae.setCompleted()
	core.Info("team detail ", detail)
}

func (r *Robot) isLeader(leaderId int64) bool {
	return r.humanId == leaderId
}