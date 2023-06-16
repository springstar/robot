package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
)

func (r *Robot) enterStage() {
	packet := msg.SerializeCSStageEnter(msg.MSG_CSStageEnter)	
	r.sendPacket(packet)
	queueMsgStat(STAT_SEND_PACKETS, int32(msg.MSG_CSStageEnter), int32(len(packet)))
}

func (r *Robot) handleEnterStage(packet *core.Packet) {
	msg := msg.ParseSCStageEnterResult(int32(packet.Type), packet.Data)
	stageObjs := msg.GetObj()
	for _, obj := range stageObjs {
		wo := newWorldObj(obj.GetObjId())
		wo.typ = WorldObjType(obj.GetType())
		wo.pos = core.NewVec2(obj.GetPos().GetX(), obj.GetPos().GetY())
		r.addVisibleObj(wo)
	}

	core.Info("enter stage, profession: ", r.profession)
	core.Info("map ", r.mapSn)
	r.fsm.trigger(r.fsm.state, "enterok", r)
	queueStat(STAT_ENTER_STAGE, 1)
}

func (r *Robot) handleSwitchStage(packet *core.Packet) {
	msg := msg.ParseSCStageSwitch(int32(packet.Type), packet.Data)
	stageId := msg.GetStageId()
	mapSn := msg.GetStageSn()
	repSn := msg.GetRepSn()
	pos := msg.GetPos()
	dir := msg.GetDir()
	lineNum := msg.GetLineNum()

	r.mapSn = mapSn
	core.Info(stageId, mapSn, repSn, pos, dir, lineNum)
	core.Info("switch stage curr state ", r.fsm.state)
	r.fsm.trigger(r.fsm.state, "switch", r)
	queueStat(STAT_SWITCH_STAGE, 1)

}

func (r *Robot) handleObjAppear(packet *core.Packet) {
	msg := msg.ParseSCStageObjectAppear(int32(msg.MSG_SCStageObjectAppear), packet.Data)
	obj := msg.GetObjAppear()
	if r.humanId == obj.GetObjId() {
		return
	}

	wo := newWorldObj(obj.GetObjId())
	wo.typ = WorldObjType(obj.GetType())
	wo.pos = core.NewVec2(obj.GetPos().GetX(), obj.GetPos().GetY())
	r.addVisibleObj(wo)
}

func (r *Robot) handleObjDisappear(packet *core.Packet) {
	msg := msg.ParseSCStageObjectDisappear(int32(msg.MSG_SCStageObjectDisappear), packet.Data)
	objId := msg.GetObjId()
	r.removeVisibleObj(objId)
}

func (r *Robot) handleStageMove(packet *core.Packet) {
	// fmt.Println("sc stage move")
}
