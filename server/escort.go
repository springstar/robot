package server

import (
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/core"
)

type EscortQuestData struct {
	questSn int
	path []*core.Vec2
	idx int

}

func newEscortQuestData(questSn int) *EscortQuestData {
	return &EscortQuestData{
		questSn: questSn,
		idx: 0,
	}
}

func (d *EscortQuestData) genPath(confQuest *config.ConfQuest) {
	target := core.Str2StrSlice(confQuest.Target)
	sn := target[1]
	confScenePath := config.FindConfScenePath(sn)
	if confScenePath == nil {
		core.Info("escort genPath no scene path found ", sn)
		return
	}

	pointSnList := core.Str2StrSlice(confScenePath.PointSNList)
	for _, p := range pointSnList {
		confScenePoint := config.FindConfScenePoint(p)
		if confScenePoint == nil {
			core.Info("escort genPath no scene point found ", p)
			continue
		}

		posSlice := core.Str2Float32Slice(confScenePoint.Position)
		pos := core.NewVec2(posSlice[0], posSlice[2])
		d.path = append(d.path, pos)
	}

}

func (d *EscortQuestData) hasNext() bool {
	return (len(d.path) > 0 && d.idx < len(d.path) - 1)
}

func (d *EscortQuestData) next() {
	d.idx += 1
}

func (d *EscortQuestData) getEscortPos() *core.Vec2 {
	if d.idx < 0 || d.idx > len(d.path) - 1 {
		return nil
	}

	return d.path[d.idx]
}

func (d *EscortQuestData) onStatusUpdate(e *RobotQuestExecutor, sn int, status QuestStatus) {
	core.Info("esocrt data onStatusUpdate ", sn, status)
	if status == QSTATE_COMPLETED {
		e.commitQuest(sn)
	}
}

func (d *EscortQuestData) dumpPath() {
	core.Info("dump escort path ")
	for _, pos := range d.path {
		core.Info("dum path pos ", pos)
	}
}

func (d *EscortQuestData) resume(e *RobotQuestExecutor) {
	core.Info("EscortQuestData resume")
	pos := d.getEscortPos()
	if pos == nil {
		d.dumpPath()
		return
	}

	ret := e.move(pos)
	if ret == -1 {
		core.Info("resume escort moving ", pos.X, pos.Y, e.pos.X, e.pos.Y)
		return
	}

	if d.hasNext() {
		d.next()
	} else {
		e.setCompleted()
		core.Info("escort finished")
	}
}

func (d *EscortQuestData) getQuestSn() int {
	return d.questSn
}