package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFsm(t *testing.T) {
	accountMgr := newAccountManager(1000, 2000)
	robotMgr := newRobotManager()
	fsm := newDefaultStateMachine(defaultTransitions())
	r := newRobot()
	r.startup(accountMgr, robotMgr, fsm)

	r.fsm.Trigger("waitForConnect", "connect", r)
	assert.Equal(t, "connected", r.state)
	r.fsm.Trigger("connected", "done", r)
	assert.Equal(t, "finished", r.state)
}