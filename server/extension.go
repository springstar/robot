package server

import (
	"github.com/springstar/robot/core"
)

type ModuleManager struct {
	modules map[string]RobotModule
}

type RobotModule interface {
	core.IModule
	core.ISubscriber

}

func newModuleManager() *ModuleManager {
	return &ModuleManager{
		modules: make(map[string]RobotModule),
	}
}

