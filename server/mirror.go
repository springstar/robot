package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
)

type StageClearQuestData struct {
	clearCount int
	enemyPosList []*core.Vec2
}

func newStageClearQuestData() *StageClearQuestData {
	return &StageClearQuestData{

	}
}

func getStageClearTarget(confQuest *config.ConfQuest) (count int, repSn int){
	infos, err := core.Str2IntSlice(confQuest.Target)
	if err != nil {
		core.Error("StageClear quest target error ", confQuest.Sn)
		return count, repSn
	}
	return infos[0], infos[1]
}

func (d *StageClearQuestData) getEnemyPosList(confQuest *config.ConfQuest) {
	c, repSn := getStageClearTarget(confQuest)
	d.clearCount = c
	confScene := config.FindConfScene(repSn)
	if confScene == nil {
		return
	}

	plotIds := core.Str2StrSlice(confScene.PlotIDs)
	for _, plotId := range plotIds {
		confScenePlot := config.FindConfScenePlot(plotId)
		if confScenePlot == nil {
			continue
		}

		evIdList := core.Str2StrSlice(confScenePlot.DoEventIDs)
		for _, eid := range evIdList {
			confSceneEvent := config.FindConfSceneEvent(eid)
			if confSceneEvent == nil {
				continue
			}

			if confSceneEvent.EventType != SEVT_FLUSHMONSTER {
				continue
			}

			if confSceneEvent.ArrfParam1 == "" {
				continue
			}

			sceneMonsterIds, err := core.Str2IntSlice(confSceneEvent.ArrfParam1)
			if err != nil {
				continue
			}

			for _, smid := range sceneMonsterIds {
				confSceneCharacter := config.FindConfSceneCharacter(smid)
				if confSceneCharacter == nil {
					continue
				}

				monsterSn, _ := core.Str2Int(confSceneCharacter.MonsterSn)
				confCharacterMonster := config.FindConfCharacterMonster(monsterSn)
				if confCharacterMonster == nil {
					continue
				}

				if confCharacterMonster.Camp == CAMP_ENEMY || confCharacterMonster.Camp == CAPM_MONSTER {
					enemyPos := core.Str2Float32Slice(confSceneCharacter.Position)

					d.enemyPosList = append(d.enemyPosList, core.NewVec2(enemyPos[0], enemyPos[2]))
				}

			}
		}
	}

}

