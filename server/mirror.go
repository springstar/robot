package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
)

type StageClearQuestData struct {
	questSn int
	clearCount int
	curEnemy int64
	snList []int
	posList []*core.Vec2
	idx int
}

func newStageClearQuestData(sn int) *StageClearQuestData {
	return &StageClearQuestData{
		questSn: sn,
		curEnemy: 0,
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

func (d *StageClearQuestData) genEnemyPosList(confQuest *config.ConfQuest) {
	c, repSn := getStageClearTarget(confQuest)
	d.clearCount = c
	confScene := config.FindConfScene(repSn)
	if confScene == nil {
		core.Info("no rep found in StageClearQuest ", repSn)
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

			sceneMonsterIds, err := core.Str2IntSlice(confSceneEvent.ArrnParam1)
			if err != nil {
				core.Info(err)
				continue
			}

			for _, smid := range sceneMonsterIds {
				confSceneCharacter := config.FindConfSceneCharacter(smid)
				if confSceneCharacter == nil {
					// core.Info("genEnemyPosList no sceneChar ", smid)
					continue
				}

				monsterSn := confSceneCharacter.MonsterSn
				confCharacterMonster := config.FindConfCharacterMonster(monsterSn)
				if confCharacterMonster == nil {
					// core.Info("genEnemyPosList no monster ", monsterSn)
					continue
				}

				if confCharacterMonster.Camp == CAMP_ENEMY || confCharacterMonster.Camp == CAPM_MONSTER {
					enemyPos := core.Str2Float32Slice(confSceneCharacter.Position)
					d.posList = append(d.posList, core.NewVec2(enemyPos[0], enemyPos[2]))
					d.snList = append(d.snList, smid)
					
				}
			}
		}
	}

}

func (d *StageClearQuestData) onStatusUpdate(e *RobotQuestExecutor, sn int, status QuestStatus) {
	core.Info("StageClearQuestData data onStatusUpdate ", sn, status)
	if status == QSTATE_COMPLETED {
		core.Info("StageClearQuest leave stage ", sn)
		e.sendLeaveInstance()
		e.setOngoing()
	}
}

func (d *StageClearQuestData) getQuestSn() int {
	return d.questSn
}

func (d *StageClearQuestData) resume(e *RobotQuestExecutor) {
	if d.curEnemy > 0 {
		vo := e.findObj(d.curEnemy)
		if vo == nil {
			d.curEnemy = 0
			return
		}

		curObj := vo.(*MonsterObj)

		if curObj.isDead() {
			d.curEnemy = 0
			return
		}

		ret := e.move(curObj.pos)
		if ret == -1 {
			return
		}

		e.fight(d.curEnemy)
		return
	}

	enemyId := d.lockEnemy(e)	
	if enemyId == 0 {
		// core.Info("stageclear quest no enemy locked")
		return
	}

	d.curEnemy = enemyId

}

func (d *StageClearQuestData) getEnemyPos() *core.Vec2 {
	if d.idx < 0 || d.idx > len(d.posList) - 1 {
		return nil
	}

	return d.posList[d.idx]
}

func (d *StageClearQuestData) getEnemySn() int {
	if d.idx < 0 || d.idx > len(d.snList) - 1 {
		return 0
	}

	return d.snList[d.idx]
}

func (d *StageClearQuestData) next() {
	d.idx = d.idx + 1
	if d.idx > len(d.snList) - 1 {
		d.idx = 0
	}
}


func (d *StageClearQuestData) lockEnemy(e *RobotQuestExecutor) int64 {
	pos := d.getEnemyPos()
	sn := d.getEnemySn()

	ret := e.move(pos)
	if ret == -1 {
		return 0
	}

	enemy := e.findMonsterObj(sn)
	if enemy == nil {
		d.next()
		return 0
	}

	if enemy.isDead() {
		d.next()
		return 0
	}

	d.next()
	core.Info("lock monster ", sn , enemy.getId())
	return enemy.getId()
}

