package server

import (
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/config"
)



type WildExecutor struct {
	*Executor
	monsterPosList []*core.Vec2
}

func newWildExecutor(r *Robot) *WildExecutor {
	return &WildExecutor{
		Executor: newExecutor(r),
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

	if int(q.mapSn) != mapSn {
		q.switchStage(1, mapSn, mapSn)
		q.setPause()	
	}

	if q.getState() == EXEC_PAUSE {	
		q.genMonsterPosList(mapSn)
		q.attachCtxFun(asyncWild, q, mapSn)
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
	n := len(q.monsterPosList)
	rnd := core.GenRandomInt(n)
	pos := q.monsterPosList[rnd]
	ret := q.move(pos)
	if ret == -1 {
		return
	}

	
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


