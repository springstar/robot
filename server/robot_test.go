package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestFsm(t *testing.T) {
	accountMgr := newAccountManager(1000, 2000)
	_, account := accountMgr.alloc()
	robotMgr := newRobotManager()
	fsm := newFsm()
	r := newRobot(account, robotMgr, fsm)
	r.startup()

	r.fsm.trigger("entry", "connect", r)
	assert.Equal(t, "connecting", r.fsm.state)
	r.fsm.trigger("connecting", "cok", r)
	assert.Equal(t, "connected", r.fsm.state)
}