package server

import (
	"strconv"
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/pb"
)

const (
	START_QUEST_SN = 1001
)

type QuestStatus int32

const (
	QSTATE_CANACCEPT = iota
	QSTATE_ONGOING 
	QSTATE_COMPLETED
	QSTATE_REWARED
	QSTATE_EXPIRED
	QSTATE_NOOPEN
	QSTATE_WAITREFRESH
	QSTATE_FAILED
)

type Quest struct {
	sn int32
	step int32
	total int32
	status int32
	typ int32
}

func newQuest() *Quest {
	return &Quest{

	}
}

type RobotQuestSet struct {
	quests map[int]*Quest
}

func newQuestSet() *RobotQuestSet{
	return &RobotQuestSet{
		quests: make(map[int]*Quest),
	}
}

func (qs *RobotQuestSet) initQuest(quests []*pb.DQuest) {
	for _, quest := range quests {
		core.Info("init quest ", quest.Sn)
		q := newQuest()		
		q.sn = quest.Sn
		q.typ = quest.Type
		q.status = quest.Status
		q.total = quest.TargetProgress
		q.step = quest.NowProgress
		qs.addQuest(q)
	}
}

func (qs *RobotQuestSet) isEmptySet() bool {
	return len(qs.quests) == 0
}

func (qs *RobotQuestSet) isPreCompleted(sn int) bool {
	confQuest := config.FindConfQuest(sn)
	if confQuest == nil {
		return false
	}

	preSn, _ := strconv.Atoi(confQuest.PreSn)
	if preSn <= 0 {
		return true
	}

	for _, v := range qs.quests {
		if v.sn == int32(preSn) && v.status == QSTATE_COMPLETED {
			return true
		}
	}

	return false
}

func (qs *RobotQuestSet) findQuestToAccept() int32 {
	return 0
}

func (qs *RobotQuestSet) findQuestToExec() int32 {
	for k, v := range qs.quests {
		core.Info("find quest to accept ", k)
		core.Info("quest status ", v.status)
		if v.status == QSTATE_ONGOING {
			return int32(k)
		}
	}

	return 0
}

func (qs *RobotQuestSet) addQuest(q *Quest) {
	if _, ok := qs.quests[int(q.sn)]; !ok {
		qs.quests[int(q.sn)] = q
	}
}

type RobotQuestExecutor struct {
	*Executor
	*RobotQuestSet
}

func newQuestExecutor(r *Robot) *RobotQuestExecutor {
	return &RobotQuestExecutor{
		Executor: newExecutor(r),
		RobotQuestSet: newQuestSet(),
	}
}

func (q *RobotQuestExecutor) acceptQuest(quest int32) ExecState {
	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) execQuest(quest int32) ExecState {
	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) exec(params []string, delta int) ExecState {
	quest := q.findQuestToExec()
	if quest > 0 {
		return q.execQuest(quest)
	}

	quest = q.findQuestToAccept()
	if quest > 0 {
		return q.acceptQuest(quest)
	}

	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) checkIfExec(params []string) bool {
	return true
}

func (q *RobotQuestExecutor) handleBreak() {

}
	 