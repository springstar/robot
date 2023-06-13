package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
	"fmt"
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
		fmt.Println(obj.GetType())
		fmt.Println(obj.GetPos())
		fmt.Println(obj.GetName())
	}

	fmt.Println("enter stage, profession: ", r.profession)
	fmt.Println("map ", r.mapSn)
	r.fsm.trigger(r.fsm.state, "enterok", r)
}

func (r *Robot) handleSwitchStage(packet *core.Packet) {
	msg := msg.ParseSCStageSwitch(int32(packet.Type), packet.Data)
	stageId := msg.GetStageId()
	mapSn := msg.GetStageSn()
	repSn := msg.GetRepSn()
	pos := msg.GetPos()
	dir := msg.GetDir()
	lineNum := msg.GetLineNum()
	fmt.Println(stageId, mapSn, repSn, pos, dir, lineNum)
	fmt.Println("switch stage curr state ", r.fsm.state)
	r.fsm.trigger(r.fsm.state, "switch", r)
	queueStat(STAT_ENTER_STAGE, 1)

}

func (r *Robot) onSwitchStage() {
	if r.profession == 0 {
		fmt.Println("awake soul after switch stage")
		r.awakeSoul()
		return
	}

	r.pc = core.GenRandomInt(serv.icount())
	fmt.Println("after switch r.pc =", r.pc)
}

func (r *Robot) handleObjAppear(packet *core.Packet) {

}

func (r *Robot) handleObjDisappear(packet *core.Packet) {

}

func (r *Robot) handleStageMove(packet *core.Packet) {
	fmt.Println("sc stage move")
}
