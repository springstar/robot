package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
	"fmt"
)

func (r *Robot) enterStage() {
	fmt.Println("enterStage")
	packet := msg.SerializeCSStageEnter(msg.MSG_CSStageEnter)	
	r.sendPacket(packet)
}

func (r *Robot) handleEnterStage(packet *core.Packet) {
	msg := msg.ParseSCStageEnterResult(int32(packet.Type), packet.Data)
	stageObjs := msg.GetObj()
	for _, obj := range stageObjs {
		fmt.Println(obj.GetType())
		fmt.Println(obj.GetPos())
		fmt.Println(obj.GetName())
	}
}

func (r *Robot) handleSwitchStage(packet *core.Packet) {
	fmt.Println("switch stage")
	msg := msg.ParseSCStageSwitch(int32(packet.Type), packet.Data)
	stageId := msg.GetStageId()
	mapSn := msg.GetStageSn()
	repSn := msg.GetRepSn()
	pos := msg.GetPos()
	dir := msg.GetDir()
	lineNum := msg.GetLineNum()
	fmt.Println(stageId, mapSn, repSn, pos, dir, lineNum)
	r.fsm.trigger(r.fsm.state, "enterok", r)

}

func (r *Robot) handleObjAppear(packet *core.Packet) {

}

func (r *Robot) handleObjDisappear(packet *core.Packet) {

}


