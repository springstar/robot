package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
)

type GatherQuestData struct {
	snList []int
	posList []*core.Vec2
	idx int
}

func newGatherQuestData(snList []int, posList []*core.Vec2) *GatherQuestData {
	return &GatherQuestData{
		snList: snList,
		posList: posList,
		idx: 0,
	}
}

func (d *GatherQuestData) hasNext() bool {
	return (len(d.posList) > 0 && d.idx < len(d.posList) - 1)
}

func (d *GatherQuestData) next() {
	if d.hasNext() {
		d.idx += 1
	}
}

func (d *GatherQuestData) getGatherPos() *core.Vec2 {
	if d.idx < 0 || d.idx > len(d.posList) - 1 {
		return nil
	}

	return d.posList[d.idx]
}

func (d *GatherQuestData) getGatherSn() int {
	if d.idx < 0 || d.idx > len(d.posList) - 1 {
		return 0
	}

	return d.snList[d.idx]

}

func (r *Robot) gatherFirst(id int64) {
	request := msg.SerializeCSGatherFirst(uint32(msg.MSG_CSGatherFirst), id, 0)
	r.sendPacket(request)
}

func (r *Robot) gatherSecond(id int64) {
	request := msg.SerializeCSGatherSecond(uint32(msg.MSG_CSGatherSecond), id)
	r.sendPacket(request)
}

func (r *Robot) HandleGatherFirst(packet *core.Packet) {
	msg := msg.ParseSCGatherFirst(int32(msg.MSG_SCGatherFirst), packet.Data)
	objId := msg.GetId()
	robotId := msg.GetHumanId()
	if r.humanId != robotId {
		return
	}

	r.gatherSecond(objId)
}

func (r *Robot) HandleGatherSecond(packet *core.Packet) {
	msg := msg.ParseSCGatherSecond(int32(msg.MSG_SCGatherSecond), packet.Data)
	robotId := msg.GetHumanId()
	if r.humanId != robotId {
		return
	}


}	