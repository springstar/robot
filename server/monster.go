package server

import (
	"github.com/springstar/robot/config"
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

func (d *MonsterQuestData) genMonsterPosList(confQuest *config.ConfQuest) {
	monsterSn, err := core.Str2Int(confQuest.ArrParam)
	if err != nil {
		core.Error(err)
		return
	}

	mapSn, err := core.Str2Int(confQuest.Target)
	if err != nil {
		core.Error(err)
		return
	}

	confs := config.GetAllConfSceneCharacter()
	for _, c := range confs {
		if mapSn != c.SceneID {
			if len(d.posList) > 0 {
				break
			} else {
				continue
			}
		}

		monster, err := core.Str2Int(c.MonsterSn)
		if err != nil {
			continue
		}

		if monster != monsterSn {
			continue
		}

		pos := core.Str2Float32Slice(c.Position)
		d.posList = append(d.posList, core.NewVec2(pos[0], pos[2]))
	}
}

