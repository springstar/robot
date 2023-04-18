package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	accountMgr := newAccountManager(1000, 2000)
	account1 := accountMgr.alloc()
	assert.Equal(t, 1000, account1.id)
	assert.Equal(t, "robot1000", account1.name)
	account2 := accountMgr.alloc()
	assert.Equal(t, 1001, account2.id)
	assert.Equal(t, "robot1001", account2.name)
	account3 := accountMgr.alloc()
	assert.Equal(t, 1002, account3.id)
	assert.Equal(t, "robot1002", account3.name)

}