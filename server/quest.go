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

type iQuestData interface {
	onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus)
}

type Quest struct {
	sn int32
	step int32
	total int32
	status int32
	typ int32
	data iQuestData
}

func newQuest() *Quest {
	return &Quest{

	}
}

func (q *Quest) attach(data iQuestData) {
	if data != nil {
		q.data = data
	}
}

func (q *Quest) detach() {
	q.data = nil
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
	q.setOngoing()
	return q.getState()
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
		q.execDialogQuest(confQuest)
	case QT_GATHER:
		q.execGatherQuest(confQuest)
	case QT_ESCORT:
		q.execEscortQuest(confQuest)
	default:
		return EXEC_NO_START	

	}

	return q.getState()
}

func (q *RobotQuestExecutor) moveToQuestPos(confQuest *config.ConfQuest) ExecState {
	mapSn, pos := getQuestPosition(confQuest)
	if int(q.mapSn) == mapSn {
		ret := q.move(pos)
		if ret == -1 {
			core.Info("exec moving to complete ", confQuest.Sn)
			q.setRepeated()
			return EXEC_REPEATED
		}	
	} else {
		q.sendEnterInstance(mapSn, confQuest.Sn)
		q.setOngoing()
		return q.getState()
	}

	q.setCompleted()
	return q.getState()

}

func (q *RobotQuestExecutor) execDialogQuest(confQuest *config.ConfQuest) ExecState {
	es := q.moveToQuestPos(confQuest)
	if es != EXEC_COMPLETED {
		return es
	}

	if confQuest.CommitType == 1 {
		msg := msg.SerializeCSCompleteQuest(uint32(msg.MSG_CSCompleteQuest), int32(confQuest.Sn))
		q.sendPacket(msg)
	}

	q.setOngoing()
	return EXEC_ONGOING
}

func getGatherInfo(confQuest *config.ConfQuest) ([]int, []*core.Vec2) {
	var sceneCharSnList []int
	infos := []string{confQuest.Target, confQuest.ArrParam, confQuest.ArrParam2}
	for _, info := range infos {
		gather, err := core.Str2IntSlice(info)
		if err != nil {
			continue
		}

		sceneCharSn := int(gather[2])
		sceneCharSnList = append(sceneCharSnList, sceneCharSn)
	}

	var gatherObjPosList []*core.Vec2

	for _, sn := range sceneCharSnList {
		confScene := config.FindConfSceneCharacter(sn)
		if confScene == nil {
			core.Warn("no ConfSceneCharacter ", sn)
			continue
		}

		position := core.Str2Float32Slice(confScene.Position)

		pos := core.NewVec2(position[0], position[2])
		gatherObjPosList = append(gatherObjPosList, pos)
	}

	return sceneCharSnList, gatherObjPosList
}

func (q *RobotQuestExecutor) execGather(d *GatherQuestData) ExecState {
	pos := d.getGatherPos()
	ret := q.move(pos)
	if ret == -1 {
		// core.Info("exec moving to complete ", confQuest.Sn)
		q.setRepeated()
		return q.getState()
	}

	sn := d.getGatherSn()
	obj := q.findGatherObj(sn)
	if obj == nil {
		core.Info("no gather obj found ", sn, pos.X, pos.Y, q.pos.X, q.pos.Y)
		q.setRepeated()
		return q.getState()
	}

	q.stepGather(obj.id)
	q.setOngoing()

	return q.getState()
}

func (q *RobotQuestExecutor) execEscortQuest(confQuest *config.ConfQuest) ExecState {
	es := q.moveToQuestPos(confQuest)
	if es != EXEC_COMPLETED {
		return q.getState()
	}

	return q.getState()
}

func (q *RobotQuestExecutor) execGatherQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		return q.getState()
	}

	if quest.data == nil {
		snList, posList := getGatherInfo(confQuest)
		if len(posList) == 0 {
			return q.getState()
		}

		qd := newGatherQuestData(snList, posList)
		quest.attach(qd)
	}

	qd := quest.data.(*GatherQuestData)

	return q.execGather(qd)
}

func (q *RobotQuestExecutor) commitQuest(sn int) {
	request := msg.SerializeCSCommitQuestNormal(uint32(msg.MSG_CSCommitQuestNormal), int32(sn))
	q.sendPacket(request)
	q.setOngoing()
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

	if quest.data != nil {
		quest.data.onStatusUpdate(q, sn, status)
	}
}


func (q *RobotQuestExecutor) handleBreak() {

}

// func getEscortInfo(confQuest *config.ConfQuest) (repSn int, pos *core.Vec2) {


// }

func getQuestPosition(confQuest *config.ConfQuest) (mapSn int, pos *core.Vec2) {
	target, _ := core.Str2IntSlice(confQuest.Target)
	switch confQuest.Type {
	case QT_DIALOG:
	case QT_ESCORT:	
		return target[1], getQuestNpcPosition(target[2])
	default:
		return target[1], core.NewZeroVec2()	
	}

	return target[1], core.NewZeroVec2()	

}

func getQuestNpcPosition(sn int) *core.Vec2{
	confSceneChar := config.FindConfSceneCharacter(sn)
	if confSceneChar == nil {
		return core.NewZeroVec2()
	}

	position := core.Str2Float32Slice(confSceneChar.Position)

	return core.NewVec2(position[0], position[2])
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