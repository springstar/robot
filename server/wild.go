package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
)



type WildExecutor struct {
	*Executor
	snList []int
	monsterPosList []*core.Vec2
	curEnemy int64
	idx int
}

func newWildExecutor(r *Robot) *WildExecutor {
	return &WildExecutor{
		Executor: newExecutor(r),
		idx: -1,
	}
}

func (q *WildExecutor) genMonsterPosList(mapSn int) {
	confScene := config.FindConfScene(mapSn)
	if confScene == nil {
		return
	}

	monsterIdList, err := core.Str2IntSlice(confScene.MonsterIDs)
	if err != nil {
		core.Info("WildExecutor monster list error ", err)
	}

	for _, monsterId := range monsterIdList {
		confChar := config.FindConfSceneCharacter(monsterId)
		if confChar == nil {
			continue
		}

		q.snList = append(q.snList, confChar.Sn)
		p := core.Str2Float32Slice(confChar.Position)
		pos := core.NewVec2(p[0], p[2])
		q.monsterPosList = append(q.monsterPosList, pos)
	}

}

func (q *WildExecutor) exec(params []string, delta int) {
	mapSn, err := core.Str2Int(params[0])
	if err != nil {
		core.Error("wild map error %s", params)
		return
	}

	if len(q.snList) == 0 {
		q.genMonsterPosList(mapSn)
	}

	if int(q.mapSn) != mapSn {
		q.switchStage(1, mapSn, mapSn)
		q.setPause()	
	}

	if q.getState() == EXEC_PAUSE {	
		q.attachCtxFun(asyncWild, q, mapSn)
	} else {
		q.execWildFight(mapSn)
	} 
}

func (q *WildExecutor) resume(params []string, delta int) {
	q.Executor.resume()
}

func asyncWild(e iExecutor, i interface{}) {
	mapSn := i.(int)
	q := e.(*WildExecutor)
	q.execWildFight(mapSn)
}

func (q *WildExecutor) execWildFight(mapSn int) {
	if q.curEnemy > 0 {
		vo := q.findObj(q.curEnemy)
		if vo == nil {
			q.curEnemy = 0
			return
		}

		curObj := vo.(*MonsterObj)
		ret := q.move(curObj.pos)
		if ret == -1 {
			return
		}

		q.fight(q.curEnemy)
		return

	}
	
	enemyId := q.lockEnemy()	
	if enemyId == 0 {
		// core.Info("monster quest no enemy locked")
		return
	}

	q.curEnemy = enemyId

}

func (q *WildExecutor) lockEnemy() int64 {
	if q.idx == -1 {
		n := len(q.monsterPosList)
		q.idx = core.GenRandomInt(n)
	}

	pos := q.monsterPosList[q.idx]
	ret := q.move(pos)
	if ret == -1 {
		// core.Info("moving to ", pos.X, pos.Y, q.pos.X, q.pos.Y)
		return 0
	}

	sn := q.snList[q.idx]
	enemy := q.findMonsterObj(sn)
	if enemy == nil {
		// core.Info("no enemy found ", sn)
		q.idx = -1
		return 0
	}

	if enemy.isDead() {
		core.Info("enemy is dead ", sn)
		q.idx = -1
		return 0
	}

	return enemy.getId()

}

func (q *WildExecutor) onEvent(k EventKey) {
	switch (k) {
	case EK_STAGE_SWITCH:
		q.onStageSwitch()
	default:
		break	
	}
}

func (q *WildExecutor) onStageSwitch() {
	if q.getState() == EXEC_PAUSE {
		q.setResume()
	}
}


