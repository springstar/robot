package server

import (
	"github.com/springstar/robot/core"
)

type MonsterQuestData struct {
	questSn int
	posList []*core.Vec2
}

func newMonsterQuestData(sn int) *MonsterQuestData {
	return &MonsterQuestData{
		questSn: sn,
	}
}

func (d *MonsterQuestData) resume(executor *RobotQuestExecutor) {

}

func (d *MonsterQuestData) getQuestSn() int {
	return d.questSn
}

func (d *MonsterQuestData) onStatusUpdate(executor *RobotQuestExecutor, sn int, status QuestStatus) {

}
