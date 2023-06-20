package server

import (
	"strconv"
	"sort"
	"github.com/springstar/robot/msg"
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

type QuestType int32

const (
	QT_LEVEL = 1
	QT_ESCORT = 74
	QT_DIALOG = 12
	QT_GATHER = 28
)

type QuestAcceptType int32 

const (
	QAT_SYSTEM = 1
	QAT_AUTO = 2
	QAT_MANUAL = 3
	QAT_UI = 4
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
		// core.Info("sn pre cheking status", sn, preSn, v.sn, v.status)
		if v.sn == int32(preSn) && (v.status == QSTATE_COMPLETED || v.status == QSTATE_REWARED) {
			return true
		}
	}

	return false
}

func (qs *RobotQuestSet) findQuest(sn int) *Quest {
	if quest, ok := qs.quests[sn]; ok {
		return quest
	}

	return nil
}

func (qs *RobotQuestSet) findQuestToAccept() int32 {
	var canAccepts []int
	for k, v := range qs.quests {
		if v.status != QSTATE_CANACCEPT {
			continue
		}

		canAccepts = append(canAccepts, int(k))		
	}

	sort.Ints(canAccepts)
	core.Info("sorted canAccepts ", canAccepts)
	if len(canAccepts) > 0 {
		quest := canAccepts[0]
		if qs.isPreCompleted(quest) {
			return int32(quest)
		}
	}

	return 0
}

func (qs *RobotQuestSet) findQuestToExec() int32 {
	for k, v := range qs.quests {
		// core.Info("find quest to exec ", k)
		// core.Info("quest status ", v.status)
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

func (q *RobotQuestExecutor) acceptQuest(quest int) ExecState {
	core.Info("accept quest ", quest)
	request := msg.SerializeCSAcceptQuest(uint32(msg.MSG_CSAcceptQuest), int32(quest))
	q.sendPacket(request)
	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) execQuest(quest int) ExecState {
	confQuest := config.FindConfQuest(int(quest))
	if confQuest == nil {
		core.Info("no such quest ", quest)
		return EXEC_COMPLETED
	}

	if !q.isPreCompleted(quest) {
		// core.Info("pre sn no completed ", quest)
		return EXEC_COMPLETED
	}

	switch confQuest.Type {
	case QT_DIALOG:
		return q.execDialogQuest(confQuest)
	case QT_GATHER:
		return q.execGatherQuest(confQuest)
	default:
		return EXEC_NO_START	

	}

	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) execDialogQuest(confQuest *config.ConfQuest) ExecState {
	pos := getQuestPosition(confQuest)
	ret := q.move(pos)
	if ret == -1 {
		core.Info("exec moving to complete ", confQuest.Sn)
		return EXEC_ONGOING
	}

	core.Info("complete quest")
	msg := msg.SerializeCSCompleteQuest(uint32(msg.MSG_CSCompleteQuest), int32(confQuest.Sn))
	q.sendPacket(msg)

	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) execGatherQuest(confQuest *config.ConfQuest) ExecState {
	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) commitQuest() {

}

func (q *RobotQuestExecutor) exec(params []string, delta int) ExecState {
	quest := q.findQuestToExec()
	if quest > 0 {
		return q.execQuest(int(quest))
	}

	quest = q.findQuestToAccept()
	if quest > 0 {
		return q.acceptQuest(int(quest))
	}

	return EXEC_COMPLETED
}

func (q *RobotQuestExecutor) updateStatus(sn int, status QuestStatus) {
	quest := q.findQuest(sn)
	if quest != nil {
		core.Info("update status ", sn, status)
		quest.status = int32(status)	
	} else {
		core.Info("new quest ", sn, status)

		quest = newQuest()
		quest.sn = int32(sn)
		quest.status = int32(status)
		q.addQuest(quest)
	}
}

func (q *RobotQuestExecutor) checkIfExec(params []string) bool {

	return true
}

func (q *RobotQuestExecutor) handleBreak() {

}

func getQuestPosition(confQuest *config.ConfQuest) *core.Vec2 {
	target, _ := core.Str2IntSlice(confQuest.Target)
	switch confQuest.Type {
	case QT_DIALOG:
		return getDialogNpcPosition(target[2])
	default:
		return core.NewZeroVec2()	
	}

	return core.NewZeroVec2()	

}

func getDialogNpcPosition(sn int) *core.Vec2{
	confSceneChar := config.FindConfSceneCharacter(sn)
	if confSceneChar == nil {
		return core.NewZeroVec2()
	}

	position := core.Str2Float32Slice(confSceneChar.Position)

	return core.NewVec2(position[0], position[1])
}

func (r *Robot) handleQuestInfo(packet *core.Packet) {
	resp := msg.ParseSCQuestInfo(int32(msg.MSG_SCQuestInfo), packet.Data)
	executor := r.findExecutor("quest").(*RobotQuestExecutor)

	quests := resp.GetQuest()
	for _, q := range quests {
		core.Info("recv quest info ", q.GetSn(), q.GetStatus())
		executor.updateStatus(int(q.GetSn()), QuestStatus(q.GetStatus()))

	}

	executor.setCompleted()
}