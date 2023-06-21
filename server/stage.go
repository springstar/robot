package server

import (
	"github.com/springstar/robot/pb"
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

		switch wo.typ {
		case WOT_PICK:
			wo.sn = int(obj.GetPick().GetStageObjectSn())
		default:
			break	
		}
		
		core.Info("recv stage enter result ", wo.id, wo.typ, wo.sn)
		
		r.addObj(wo)
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
	// if r.humanId == obj.GetObjId() {
	// 	return
	// }

	wo := newWorldObj(obj.GetObjId())
	wo.typ = WorldObjType(obj.GetType())
	wo.pos = core.NewVec2(obj.GetPos().GetX(), obj.GetPos().GetY())
	switch wo.typ {
	case WOT_PICK:
		wo.sn = int(msg.GetObjAppear().Pick.GetStageObjectSn())
	default:
		break	
	}

	core.Info("recv obj appear ", wo.id, wo.typ, wo.sn)

	r.addObj(wo)
}

func (r *Robot) handleObjDisappear(packet *core.Packet) {
	msg := msg.ParseSCStageObjectDisappear(int32(msg.MSG_SCStageObjectDisappear), packet.Data)
	objId := msg.GetObjId()
	r.removeObj(objId)
}

func (r *Robot) handleStageMove(packet *core.Packet) {
	// msg := msg.ParseSCStageMove(int32(msg.MSG_SCStageMove), packet.Data)
	
}

func (r *Robot) sendMoveStop(x, y float32) {
	dpos := &pb.DVector2{}
	dpos.X = x
	dpos.Y = y
	ddir := &pb.DVector2{}
	ddir.X = r.dir.X
	ddir.Y = r.dir.Y
	request := msg.SerializeCSStageMoveStop(uint32(msg.MSG_CSStageMoveStop), r.humanId, dpos, ddir)
	r.sendPacket(request)
}

func (r *Robot) handleMoveStop(packet *core.Packet) {
	msg := msg.ParseSCStageMoveStop(int32(msg.MSG_SCStageMoveStop), packet.Data)
	r.pos.X = msg.GetPosEnd().GetX()
	r.pos.Y = msg.GetPosEnd().GetY()
	r.dir.X = msg.GetDir().GetX()
	r.dir.Y = msg.GetDir().GetY()
}
