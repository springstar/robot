package server

import (
	"strings"
	"strconv"
)

type Account struct {
	id   int
	name string
}

type AccountManager struct {
	startId int
	idx int
	max int
}

func newAccountManager(start int, max int) *AccountManager {
	return &AccountManager {
		startId : start,
		idx : 0,
		max : max,
	}
}

func (mgr *AccountManager) alloc() *Account{
	if mgr.idx >= mgr.max {
		return nil
	}

	account := &Account{}
	account.id = mgr.startId + mgr.idx
	var sb strings.Builder
	sb.WriteString(ROBOT_PREFIX)
	sb.WriteString(strconv.Itoa(account.id))
	account.name = sb.String()
	mgr.idx += 1
	return account
}