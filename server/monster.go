package server

import (
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/core"
)

type MonsterQuestData struct {
	questSn int
	count int
	snList []int
	posList []*core.Vec2
	curEnemy int64
	idx int
}

func newMonsterQuestData(sn int) *MonsterQuestData {
	return &MonsterQuestData{
		questSn: sn,
	}
}

func (d *MonsterQuestData) resume(e *RobotQuestExecutor) {
	if d.curEnemy > 0 {
		vo := e.findObj(d.curEnemy)
		if vo == nil {
			d.curEnemy = 0
			return
		}

		curObj := vo.(*MonsterObj)

		ret := e.move(curObj.pos)
		if ret == -1 {
			return
		}

		e.fight(d.curEnemy)
		return
	}

	enemyId := d.lockEnemy(e)	
	if enemyId == 0 {
		// core.Info("monster quest no enemy locked")
		return
	}

	d.curEnemy = enemyId

}

func (d *MonsterQuestData) getEnemyPos() *core.Vec2 {
	if d.idx < 0 || d.idx > len(d.posList) - 1 {
		return nil
	}

	return d.posList[d.idx]
}

func (d *MonsterQuestData) getEnemySn() int {
	if d.idx < 0 || d.idx > len(d.snList) - 1 {
		return 0
	}

	return d.snList[d.idx]
}

func (d *MonsterQuestData) next() {
	d.idx = d.idx + 1
	if d.idx > len(d.snList) - 1 {
		d.idx = 0
	}
}

func (d *MonsterQuestData) lockEnemy(e *RobotQuestExecutor) int64 {
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
	// core.Info("lock monster ", sn , enemy.getId())
	return enemy.getId()

}

func (d *MonsterQuestData) getQuestSn() int {
	return d.questSn
}

func (d *MonsterQuestData) onStatusUpdate(e *RobotQuestExecutor, sn int, status QuestStatus) {
	// core.Info("MonsterQuestData data onStatusUpdate ", sn, status)
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

	mapSn, err := core.Str2Int(confQuest.Target)
	if err != nil {
		core.Error(err)
		return
	}

	confs := config.GetAllConfSceneCharacter()
	for _, c := range confs {
		if mapSn != c.SceneID {
			if len(d.snList) == d.count {
				break
			} else {
				continue
			}
		}

		if !c.CanBeAttacked {
			continue
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
		d.snList = append(d.snList, c.Sn)
	}
}

