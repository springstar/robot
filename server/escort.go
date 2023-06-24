package server

import (
	"github.com/springstar/robot/config"
	"github.com/springstar/robot/core"
)

type EscortQuestData struct {
	path []*core.Vec2
}

func newEscortQuestData() *EscortQuestData {
	return &EscortQuestData{}
}

func (d *EscortQuestData) genPath(sn int) {
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

