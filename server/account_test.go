package server

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestAlloc(t *testing.T) {
	accountMgr := newAccountManager(1000, 2000)
	_, account1 := accountMgr.alloc()
	assert.Equal(t, 1000, account1.id)
	assert.Equal(t, "robot1000", account1.name)
	_, account2 := accountMgr.alloc()
	assert.Equal(t, 1001, account2.id)
	assert.Equal(t, "robot1001", account2.name)
	_, account3 := accountMgr.alloc()
	assert.Equal(t, 1002, account3.id)
	assert.Equal(t, "robot1002", account3.name)

	acc := accountMgr.findAccountByName("robot1001")
	assert.Equal(t, acc.id, account2.id)

	err, acc1 := accountMgr.findAccount(1001)
	assert.Equal(t, err, nil)
	assert.Equal(t, acc, acc1)

}