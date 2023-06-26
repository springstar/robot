package server

import (
	"github.com/springstar/robot/core"
)

func (r *Robot) handleSkillChange(packet *core.Packet) {
	core.Info("recv skill change")
}

func (r *Robot) handleHpChange(packet *core.Packet) {
	core.Info("recv hp change")
}