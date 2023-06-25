package server

import (
	"github.com/springstar/robot/msg"
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
)

type GatherQuestData struct {
	questSn int
	snList []int
	posList []*core.Vec2
	idx int
}

func newGatherQuestData(questSn int, snList []int, posList []*core.Vec2) *GatherQuestData {
	return &GatherQuestData{
		questSn: questSn,
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

func (d *GatherQuestData) resume(executor *RobotQuestExecutor) {

}

func (d *GatherQuestData) getQuestSn() int {
	return d.questSn
}

func (d *GatherQuestData) onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus) {
	if status == QSTATE_ONGOING {
		d.next()
	} else if status == QSTATE_COMPLETED {
		executor.commitQuest(sn)
	}
}

func (r *Robot) stepGather(id int64) {
	r.gatherFirst(id)
	r.gatherSecond(id)
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

	core.Info("recv gather first")


	r.gatherSecond(objId)
}

func (r *Robot) HandleGatherSecond(packet *core.Packet) {
	msg := msg.ParseSCGatherSecond(int32(msg.MSG_SCGatherSecond), packet.Data)
	robotId := msg.GetHumanId()
	if r.humanId != robotId {
		return
	}

	core.Info("recv gather second")
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