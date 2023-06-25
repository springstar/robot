package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

type TeamOperation int32

const (
	TO_NONE = iota
	TO_CREATE
	TO_LIST
)

type TeamExecutor struct {
	*Executor
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
		Executor: newExecutor(r),
		TeamMember: newTeamMember(0, false),
	}
}

func (t *TeamExecutor) resume(params []string, delta int) {

}

func (t *TeamExecutor) exec(params []string, delta int) {
	op, _ := parseOperation(params)
	switch op {
	case TO_CREATE:
		t.create(params)
	case TO_LIST:
		t.list(params)
	}
}

func (t *TeamExecutor) onEvent(k EventKey) {

}

func (t *TeamExecutor) create(params []string) ExecState {
	core.Info("send create team request")

	targetSn, _ := core.Str2Int(params[1])
	teamType, _ := core.Str2Int(params[2])
	request := msg.SerializeCSPlatCreateTeam(msg.MSG_CSPlatCreateTeam, int32(targetSn), int32(teamType))
	t.sendPacket(request)
	t.setOngoing()
	return t.getState()
}

func (t *TeamExecutor) list(params []string) ExecState {
	core.Info("request team platform list")
	targetSn, _ := core.Str2Int(params[1])
	request := msg.SerializeCSPlatTeamListRequest(msg.MSG_CSPlatTeamListRequest, int32(targetSn))
	t.sendPacket(request)
	t.setOngoing()
	return t.getState()
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
		return t.teamId <= 0	
	default:
		return true	
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
	e.setCompleted()
	core.Info("team detail ", detail)
}

func (r *Robot) isLeader(leaderId int64) bool {
	return r.humanId == leaderId
}

func (r *Robot) handleTeamList(packet *core.Packet) {
	core.Info("recv team platform list")
	resp := msg.ParseSCPlatTeamListResponse(int32(msg.MSG_SCPlatTeamListResponse), packet.Data)
	teams := resp.GetTeams()
	for _, t := range teams {
		core.Info(t)
	}

	e := r.findExecutor("team").(*TeamExecutor)
	e.setCompleted()
}