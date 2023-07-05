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
		var vo iVisible

		id := obj.GetObjId()
		typ := WorldObjType(obj.GetType())
		pos := core.NewVec2(obj.GetPos().GetX(), obj.GetPos().GetY())

		switch typ {
		case WOT_PICK:
			vo = createGather(id, typ, pos, int(obj.GetPick().GetStageObjectSn()))
		case WOT_MONSTER:
			vo = createMonster(id, typ, pos, int(obj.GetMonster().GetStageObjectSn()), obj.GetMonster().GetHpCur(), obj.GetMonster().GetHpMax())	
		default:
			vo = newWorldObj(id, typ, pos)	
		}
		
		// core.Info("recv stage enter result ", wo.id, wo.typ, wo.sn)
		
		r.addObj(vo)
	}

	if r.fsm.state == "enterstage" {
		r.fireEvent(EK_STAGE_SWITCH)
	}

	core.Info("enter stage, map sn: ", r.name, r.mapSn)
	r.fsm.trigger(r.fsm.state, "enterok", r)
	queueStat(STAT_ENTER_STAGE, 1)

	r.sendGM("-gm buff addBuff 9999999")

}

func (r *Robot) handleSwitchStage(packet *core.Packet) {
	msg := msg.ParseSCStageSwitch(int32(packet.Type), packet.Data)
	// stageId := msg.GetStageId()
	mapSn := msg.GetStageSn()
	// repSn := msg.GetRepSn()
	pos := msg.GetPos()
	dir := msg.GetDir()
	// lineNum := msg.GetLineNum()

	r.mapSn = mapSn
	r.pos.Set(pos.GetX(), pos.GetY())
	r.dir.Set(dir.GetX(), dir.GetY())
	// core.Info(stageId, mapSn, repSn, pos, dir, lineNum)
	// core.Info("switch stage curr state ", r.fsm.state)
	r.fsm.trigger(r.fsm.state, "switch", r)
	queueStat(STAT_SWITCH_STAGE, 1)

}

func (r *Robot) handleObjAppear(packet *core.Packet) {
	msg := msg.ParseSCStageObjectAppear(int32(msg.MSG_SCStageObjectAppear), packet.Data)
	obj := msg.GetObjAppear()
	// if r.humanId == obj.GetObjId() {
	// 	return
	// }

	id := obj.GetObjId()
	typ := WorldObjType(obj.GetType())
	pos := core.NewVec2(obj.GetPos().GetX(), obj.GetPos().GetY())
	var vo iVisible	
	switch typ {
	case WOT_PICK:
		sn := int(msg.GetObjAppear().Pick.GetStageObjectSn())
		vo = createGather(id, typ, pos, sn)
	case WOT_MONSTER:
		sn := int(msg.GetObjAppear().Monster.GetStageObjectSn())	
		vo = createMonster(id, typ, pos, sn, obj.GetMonster().GetHpCur(), obj.GetMonster().GetHpMax())
		// mo := vo.(*MonsterObj)
		// core.Info("recv obj appear ", vo.getId(), vo.getType(), mo.sn, r.pos.X, r.pos.Y)
	default:
		vo = newWorldObj(id, typ, pos)	
	}



	r.addObj(vo)
}

func (r *Robot) handleObjDisappear(packet *core.Packet) {
	msg := msg.ParseSCStageObjectDisappear(int32(msg.MSG_SCStageObjectDisappear), packet.Data)
	objId := msg.GetObjId()
	// core.Info("recv obj disappear ", objId, r.pos.X, r.pos.Y)
	r.removeObj(objId)
}

func (r *Robot) handleStageMove(packet *core.Packet) {
	msg := msg.ParseSCStageMove(int32(msg.MSG_SCStageMove), packet.Data)
	objId := msg.GetObjId()
	obj := r.findObj(objId)
	if obj == nil {
		return
	}

	posBegin := msg.GetPosBegin()
	posEnd := msg.GetPosEnd()
	if len(posEnd) == 0 {
		obj.setPos(posBegin.GetX(), posBegin.GetY())
		return
	}

	for _, p := range posEnd {
		obj.setPos(p.GetX(), p.GetY())
	}
	
}

func (r *Robot) handleStageSetPos(packet *core.Packet) {
	msg := msg.ParseSCStageSetPos(int32(msg.MSG_SCStageSetPos), packet.Data)
	humanId := msg.GetId()
	if r.humanId != humanId {
		return
	}

	pos := msg.GetPos()
	dir := msg.GetDir()
	r.pos.Set(pos.GetX(), pos.GetY())
	r.dir.Set(dir.GetX(), dir.GetY())
	
}

func (r *Robot) handleStagePull(packet *core.Packet) {
	msg := msg.ParseSCStagePullTo(int32(msg.MSG_SCStagePullTo), packet.Data)
	pos := msg.GetPos()
	dir := msg.GetDir()
	r.pos.Set(pos.GetX(), pos.GetY())
	r.dir.Set(dir.GetX(), dir.GetY())
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
