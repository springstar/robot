package server

import (
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/core"
)

type MonsterQuestData struct {
	questSn int
	count int
	monsterSn int
	monsters map[int]*core.Vec2
	curEnemy int64
}

func newMonsterQuestData(sn int) *MonsterQuestData {
	return &MonsterQuestData{
		questSn: sn,
		monsters: make(map[int]*core.Vec2),
	}
}

func (d *MonsterQuestData) resume(e *RobotQuestExecutor) {
	if d.curEnemy > 0 {
		curObj := e.findObj(d.curEnemy).(*MonsterObj)
		if curObj.isDead() {
			d.curEnemy = 0
			return
		}
		e.fight(d.curEnemy)
		return
	}

	enemyId := d.lockEnemy(e)	
	if enemyId == 0 {
		core.Info("monster quest no enemy locked")
		return
	}

	d.curEnemy = enemyId
	e.fight(enemyId)

}

func (d *MonsterQuestData) lockEnemy(e *RobotQuestExecutor) int64 {
	for sn, pos := range d.monsters {
		core.Info("MonsterQuestData locking enemy ", sn)
		ret := e.move(pos)
		if ret == -1 {
			continue
		}

		enemy := e.findMonsterObj(sn)
		if enemy == nil {
			core.Info("monster quest no monster ", sn)
			continue
		}

		if enemy.isDead() {
			core.Info("monster already dead ", sn)
			continue
		}

		core.Info("lock monster ", sn , enemy.getId())
		return enemy.getId()

	}

	return 0
}

func (d *MonsterQuestData) getQuestSn() int {
	return d.questSn
}

func (d *MonsterQuestData) onStatusUpdate(e *RobotQuestExecutor, sn int, status QuestStatus) {
	core.Info("MonsterQuestData data onStatusUpdate ", sn, status)
	if status == QSTATE_COMPLETED {
		core.Info("MonsterQuestData commit quest ", sn)
		e.commitQuest(sn)
	}
}

func (d *MonsterQuestData) genMonsterInfo(confQuest *config.ConfQuest) {
	count, err := core.Str2Int(confQuest.ArrParam2)
	if err != nil {
		core.Error(err)
		return
	}

	d.count = count

	monsterSn, err := core.Str2Int(confQuest.ArrParam)
	if err != nil {
		core.Error(err)
		return
	}

	d.monsterSn = monsterSn

	mapSn, err := core.Str2Int(confQuest.Target)
	if err != nil {
		core.Error(err)
		return
	}

	confs := config.GetAllConfSceneCharacter()
	for _, c := range confs {
		if mapSn != c.SceneID {
			if len(d.monsters) == d.count {
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
		d.monsters[c.Sn] = core.NewVec2(pos[0], pos[2])
		
	}
}

