package server

import (
	_ "fmt"
	"strings"
	"strconv"
	"errors"
	_ "github.com/springstar/robot/core"
	_ "github.com/springstar/robot/msg"
)

var (
	ErrExceedMaximum = errors.New("account exceed ")
)

type Account struct {
	id   int
	name string
}

type AccountManager struct {
	startId int
	idx int
	max int
	accounts map[int]*Account
}

func newAccountManager(start int, max int) *AccountManager {
	return &AccountManager {
		startId : start,
		idx : 0,
		max : max,
		accounts : make(map[int]*Account),
	}
}

func (m *AccountManager) init(r *Robot) {
	
}

func (m *AccountManager) addAccount(aid int, account *Account) {
	if _, ok := m.accounts[aid]; !ok {
		m.accounts[aid] = account
	}
}

func (m *AccountManager) findAccount(aid int) (error, *Account) {
	if _, ok := m.accounts[aid]; ok {
		return nil, m.accounts[aid]
	}

	return errors.New("no account found"), nil
}

func (m *AccountManager) findAccountByName(name string) *Account {
	for _, account := range m.accounts {
		if account.name == name {
			return account
		}
	}

	return nil
}

func (m *AccountManager) alloc() (error, *Account) {
	if m.idx >= m.max {
		return ErrExceedMaximum, nil
	}

	account := &Account{}
	account.id = m.startId + m.idx
	var sb strings.Builder
	sb.WriteString(ROBOT_PREFIX)
	sb.WriteString(strconv.Itoa(account.id))
	account.name = sb.String()
	m.idx += 1
	m.addAccount(account.id, account)

	return nil, account
}

func (m *AccountManager) allocName(name string) (error, *Account) {
	if m.idx >= m.max {
		return ErrExceedMaximum, nil
	}

	account := &Account{}
	account.id = m.startId + m.idx
	account.name = name
	m.idx += 1	
	m.addAccount(account.id, account)

	return nil, account

}

