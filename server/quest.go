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
	QT_SKILL = 5
	QT_DIALOG = 12
	QT_STAGECLEAR = 13
	QT_EXPLORE = 28
	QT_DELIVERY = 71
	QT_WAIT = 72
	QT_GATHER = 73
	QT_ESCORT = 74
	QT_TRANSFORM = 75
	QT_SOUL = 77
	QT_MONSTER = 81

)

type QuestAcceptType int32 

const (
	QAT_SYSTEM = 1
	QAT_AUTO = 2
	QAT_MANUAL = 3
	QAT_UI = 4
)

type iQuestData interface {
	getQuestSn() int
	onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus)
	resume(executor *RobotQuestExecutor)
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
		data: nil,
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
	// 临时处理，忽略非主线任务
	for _, quest := range quests {
		if quest.Status == 0 {
			continue
		}

		confQuest := config.FindConfQuest(int(quest.Sn))
		if confQuest.Group != 1 {
			continue
		}

		core.Info("init quest ", quest.Sn, quest.Status)

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

	quest := qs.findQuest(sn)
	if quest != nil && quest.status == QSTATE_ONGOING {
		return true
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
	// core.Info("sorted canAccepts ", canAccepts)
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

func (qs *RobotQuestSet) removeQuest(questSn int) {
	if _, ok := qs.quests[questSn]; ok {
		delete(qs.quests, questSn)
	}
}

type RobotQuestExecutor struct {
	*Executor
	*RobotQuestSet
	curQuest int
}

func newQuestExecutor(r *Robot) *RobotQuestExecutor {
	return &RobotQuestExecutor{
		Executor: newExecutor(r),
		RobotQuestSet: newQuestSet(),
	}
}

func (q *RobotQuestExecutor) acceptQuest(quest int) ExecState {
	// core.Info("accept quest ", quest)
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
		core.Info("pre sn no completed ", quest)
		return EXEC_COMPLETED
	}

	q.curQuest = quest

	switch confQuest.Type {
	case QT_DIALOG:
		q.execDialogQuest(confQuest)
	case QT_EXPLORE:
		q.execGatherQuest(confQuest)
	case QT_GATHER:
		q.execGatherQuest(confQuest)	
	case QT_ESCORT:
		q.execEscortQuest(confQuest)
	case QT_STAGECLEAR:
		q.execStageClearQuest(confQuest)
	case QT_SKILL:
		q.execSkillQuest(confQuest)
	case QT_MONSTER:
		q.execMonsterQuest(confQuest)
	case QT_DELIVERY:
		q.execDeliverQuest(confQuest)
	case QT_SOUL:
		q.execSoulQuest(confQuest)
	case QT_TRANSFORM:
		q.execTransformQuest(confQuest)
	case QT_WAIT:
		q.execWaitQuest(confQuest)	
	default:	
		break	

	}

	return q.getState()
}

func (q *RobotQuestExecutor) moveToQuestPos(confQuest *config.ConfQuest) {
	if q.getState() == EXEC_ONGOING {
		return
	}

	mapSn, pos := getQuestPosition(confQuest)
	if int(q.mapSn) == mapSn {
		ret := q.move(pos)
		if ret == -1 {
			// core.Info("exec moving to complete ", confQuest.Sn, mapSn, pos.X, pos.Y)
			q.setRepeated()
		} else {
			// core.Info("move ok ", confQuest.Sn, mapSn, pos, q.getState())
			q.setCompleted()
		}	
		return
	}

	// core.Info("enter instarnce ", mapSn, confQuest.Sn)
	q.sendEnterInstance(mapSn, confQuest.Sn)
	q.setPause()
}

func (q *RobotQuestExecutor) execDialogQuest(confQuest *config.ConfQuest) ExecState {
	q.moveToQuestPos(confQuest)
	if q.getState() != EXEC_COMPLETED {
		return q.getState()
	}

	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		q.setCompleted()
		return q.getState()
	}

	if quest.data == nil {
		qd := newSkillQuestData()
		quest.attach(qd)
	}

	if quest.status == QSTATE_ONGOING {
		q.completeQuest(confQuest.Sn)
	}

	// core.Info("dialog quest ongoing ", confQuest.Sn)

	q.setRepeated()
	
	return q.getState()
}

func (q *RobotQuestExecutor) execDeliverQuest(confQuest *config.ConfQuest) ExecState {
	q.moveToQuestPos(confQuest)
	if q.getState() != EXEC_COMPLETED {
		return q.getState()
	}

	if confQuest.CommitType == 1 {
		msg := msg.SerializeCSCompleteQuest(uint32(msg.MSG_CSCompleteQuest), int32(confQuest.Sn))
		q.sendPacket(msg)
	}

	q.setOngoing()

	return q.getState()
}

func (q *RobotQuestExecutor) execWaitQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		q.setCompleted()
		return q.getState()
	}

	if quest.data == nil {
		qd := newSkillQuestData()
		quest.attach(qd)
	}
	
	return q.getState()
}

func (q *RobotQuestExecutor) execSoulQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		q.setCompleted()
		return q.getState()
	}
	// 先用SkillQuestData，其实应该是一个DumbAsyncQuestData
	if quest.data == nil {
		qd := newSkillQuestData()
		quest.attach(qd)
	}

	q.sendUpgradeSoulRequest()
	q.setOngoing()
	return q.getState()
}

func (q *RobotQuestExecutor) execSkillQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		q.setCompleted()
		return q.getState()
	}

	if quest.data == nil {
		qd := newSkillQuestData()
		quest.attach(qd)
	}

	q.upgradeSkill()
	q.setOngoing()
	return q.getState()
}

func (q *RobotQuestExecutor) execTransformQuest(confQuest *config.ConfQuest) {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		q.setCompleted()
		return
	}

	q.moveToQuestPos(confQuest)
	if q.getState() == EXEC_REPEATED {
		return
	}

	//借用一下SkillQuestData，以后改掉
	if quest.data == nil {
		qd := newSkillQuestData()
		quest.attach(qd)
	}

	q.completeQuest(confQuest.Sn)
}

func (q *RobotQuestExecutor) execGather(d *GatherQuestData)  {
	pos := d.getGatherPos()
	ret := q.move(pos)
	if ret == -1 {
		q.setRepeated()
	}

	sn := d.getGatherSn()
	obj := q.findGatherObj(sn)
	if obj == nil {
		q.setRepeated()
		return
	}

	q.stepGather(obj.getId())
	q.setRepeated()

}

func (q *RobotQuestExecutor) execMonsterQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		return q.getState()
	}

	if quest.data == nil {
		// core.Info("attach kill monster quest data ", confQuest.Sn, q.getState(), q.mapSn)
		qd := newMonsterQuestData(confQuest.Sn)
		qd.genMonsterInfo(confQuest)
		quest.attach(qd)
	} 
	
	q.moveToQuestPos(confQuest)
	if q.getState() == EXEC_PAUSE {
		// core.Info("monster quest attach ctx function ", q.mapSn)
		q.attachCtxFun(asyncKillMonster, q, nil)
	} else {
		quest.data.resume(q)
	}

	return q.getState()
}

func (q *RobotQuestExecutor) execStageClearQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		return q.getState()
	}

	if quest.data == nil {
		// core.Info("attach stage clear quest data ", confQuest.Sn, q.getState())
		qd := newStageClearQuestData(confQuest.Sn)
		qd.genEnemyPosList(confQuest)
		quest.attach(qd)
	} 

	q.moveToQuestPos(confQuest)
	if q.getState() == EXEC_PAUSE {	
		q.attachCtxFun(asyncStageClear, q, nil)
	} else {
		quest.data.resume(q)
	}

	return q.getState()
	
}

func (q *RobotQuestExecutor) execEscortQuest(confQuest *config.ConfQuest) ExecState {
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		return q.getState()
	}

	if quest.data == nil {
		core.Info("attach escort quest data ", confQuest.Sn)
		qd := newEscortQuestData(confQuest.Sn)
		qd.genPath(confQuest)
		quest.attach(qd)
	} else {
		// core.Info("no need to attach escort quest data ", confQuest.Sn)
	}

	q.moveToQuestPos(confQuest)
	if q.getState() == EXEC_PAUSE {
		q.attachCtxFun(asyncEscort, q, nil)
	}

	return q.getState()
}

func asyncKillMonster(e iExecutor, i interface{}) {
	q := e.(*RobotQuestExecutor)
	q.execKillMonster()
}

// todo : refactor duplicate function
func (q *RobotQuestExecutor) execKillMonster() {
	quest := q.findQuest(q.curQuest)
	if quest == nil {
		core.Info("kill monster no quest found ", q.curQuest)
		return
	}

	if quest.data != nil {
		quest.data.resume(q)
	} else {
		core.Info("kill monster quest data nil ", q.curQuest)
	}
}

func asyncStageClear(e iExecutor, i interface{}) {
	q := e.(*RobotQuestExecutor)
	q.execStageClear()
}

func (q *RobotQuestExecutor) execStageClear() {
	quest := q.findQuest(q.curQuest)
	if quest == nil {
		core.Info("execStageClear no quest found ", q.curQuest)
		return
	}

	if quest.data != nil {
		quest.data.resume(q)
	} else {
		core.Info("stage clear quest data nil ", q.curQuest)
	}
}

func asyncEscort(e iExecutor, i interface{}) {
	q := e.(*RobotQuestExecutor)
	q.execEscort()
}

func (q *RobotQuestExecutor) execEscort() {
	quest := q.findQuest(q.curQuest)
	if quest == nil {
		core.Info("execEscort no quest found ", q.curQuest)
		return
	}

	if quest.data != nil {
		quest.data.resume(q)
	} else {
		core.Info("escort quest data nil ", q.curQuest)
	}
}

func (q *RobotQuestExecutor) execGatherQuest(confQuest *config.ConfQuest)  {
	// core.Info("exec gather quest ", confQuest.Sn)
	quest := q.findQuest(confQuest.Sn)
	if quest == nil {
		return
	}

	if quest.data == nil {
		// core.Info("attach gather quest data ", confQuest.Sn)
		qd := newGatherQuestData(confQuest.Sn)
		qd.genGatherInfo(confQuest)
		quest.attach(qd)
	}

	qd := quest.data.(*GatherQuestData)

	q.execGather(qd)
}

func (q *RobotQuestExecutor) completeQuest(sn int) {
	confQuest := config.FindConfQuest(sn)
	if confQuest == nil {
		return
	}

	msg := msg.SerializeCSCompleteQuest(uint32(msg.MSG_CSCompleteQuest), int32(confQuest.Sn))
	q.sendPacket(msg)
}

func (q *RobotQuestExecutor) commitQuest(sn int) {
	confQuest := config.FindConfQuest(sn)
	if confQuest == nil {
		return
	}

	if confQuest.CommitType == 1 {
		return
	}

	request := msg.SerializeCSCommitQuestNormal(uint32(msg.MSG_CSCommitQuestNormal), int32(sn))
	q.sendPacket(request)
	// q.setOngoing()
}



func (q *RobotQuestExecutor) exec(params []string, delta int) {
	// core.Info("RobotQuestExecutor exec")
	quest := q.findQuestToExec()
	if quest > 0 {
		q.execQuest(int(quest))
		return
	}

	quest = q.findQuestToAccept()
	if quest > 0 {
		q.acceptQuest(int(quest))
	}
}

func (q *RobotQuestExecutor) resume(params []string, delta int) {
	// core.Info("RobotQuestExecutor resume exec")
	q.Executor.resume()
}

func (q *RobotQuestExecutor) onEvent(k EventKey) {
	switch (k) {
	case EK_STAGE_SWITCH:
		q.onStageSwitch()
	default:
		break	
	}
}

func (q *RobotQuestExecutor) onStageSwitch() {
	if q.getState() == EXEC_PAUSE {
		// core.Info("set resume state")
		q.setResume()
	}

	quest := q.findQuest(q.curQuest)
	if quest == nil {
		core.Info("find quest to commit quest after stage switch ", q.curQuest)
		return
	}

	if quest.status == QSTATE_COMPLETED {
		q.commitQuest(q.curQuest)
	}
}

func (q *RobotQuestExecutor) updateStatus(sn int, status QuestStatus) {
	quest := q.findQuest(sn)
	if quest != nil {
		// core.Info("update status ", sn, status)
		quest.status = int32(status)	
	} else {
		core.Info("new quest ", sn, status, q.getState())

		quest = newQuest()
		quest.sn = int32(sn)
		quest.status = int32(status)
		q.addQuest(quest)

		q.setCompleted()
	}

	if quest.data != nil {
		quest.data.onStatusUpdate(q, sn, status)
	}
}

func getQuestPosition(confQuest *config.ConfQuest) (mapSn int, pos *core.Vec2) {
	target, _ := core.Str2IntSlice(confQuest.Target)
	switch confQuest.Type {
	case QT_DIALOG:
		return target[1], getQuestNpcPosition(target[2])
	case QT_ESCORT:
		return target[1], nil	
	case QT_STAGECLEAR:
		return target[1], nil
	case QT_MONSTER:
		return target[0], nil	
	case QT_DELIVERY:
		return getQuestNpcMap(target[2]), getQuestNpcPosition(target[2])
	case QT_TRANSFORM:
		return getQuestNpcMap(target[1]), getQuestNpcPosition(target[1])
	}

	return target[1], core.NewZeroVec2()	

}

func getQuestNpcMap(sn int) int {
	confSceneChar := config.FindConfSceneCharacter(sn)
	if confSceneChar == nil {
		return 0
	}

	return confSceneChar.SceneID
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
		confQuest := config.FindConfQuest(int(q.GetSn())) 
		if confQuest.Group != 1 {
			continue
		}

		// core.Info("recv quest info ", q.GetSn(), q.GetStatus())
		executor.updateStatus(int(q.GetSn()), QuestStatus(q.GetStatus()))
	}

	executor.setCompleted()
}

func (r *Robot) handleRemoveQuest(packet *core.Packet) {
	msg := msg.ParseSCRemoveQuest(int32(msg.MSG_SCRemoveQuest), packet.Data)
	executor := r.findExecutor("quest").(*RobotQuestExecutor)

	quests := msg.GetQuestIds()
	for _, q := range quests {
		executor.removeQuest(int(q))
	}
}