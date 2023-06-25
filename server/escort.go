package server

import (
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/core"
)

type EscortQuestData struct {
	questSn int
	path []*core.Vec2
}

func newEscortQuestData(questSn int) *EscortQuestData {
	return &EscortQuestData{
		questSn: questSn,
	}
}

func (d *EscortQuestData) genPath(confQuest *config.ConfQuest) {
	target, _ := core.Str2IntSlice(confQuest.Target)
	sn := target[2]
	confScenePath := config.FindConfScenePath(sn)
	if confScenePath == nil {
		return
	}

	points := confScenePath.PointSNList
	for _, p := range points {
		confScenePoint := config.FindConfScenePoint(p)
		if confScenePoint == nil {
			continue
		}

		posSlice := core.Str2Float32Slice(confScenePoint.Position)
		pos := core.NewVec2(posSlice[0], posSlice[2])
		d.path = append(d.path, pos)
	}

}

func (d *EscortQuestData) onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus) {

}

func (d *EscortQuestData) resume(executor *RobotQuestExecutor) {
	
}

func (d *EscortQuestData) getQuestSn() int {
	return d.questSn
}