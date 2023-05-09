package server

import (
	"fmt"
	"strings"
	"strconv"
	"errors"
	"github.com/springstar/robot/core"
	"github.com/springstar/robot/msg"
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
}

func newAccountManager(start int, max int) *AccountManager {
	return &AccountManager {
		startId : start,
		idx : 0,
		max : max,
	}
}

func (mgr *AccountManager) init(r *Robot) {
	r.dispatcher.Register(112, mgr)
}

func (mgr *AccountManager) HandleMessage(packet *core.Packet) {
	fmt.Println("recv packet ", packet.Type)
	msg := msg.ParseSCLoginResult(int32(packet.Type), packet.Data)
	fmt.Println(msg.GetResultCode())
	fmt.Println(msg.GetResultReason())
}

func (mgr *AccountManager) alloc() (error, *Account) {
	if mgr.idx >= mgr.max {
		return ErrExceedMaximum, nil
	}

	account := &Account{}
	account.id = mgr.startId + mgr.idx
	var sb strings.Builder
	sb.WriteString(ROBOT_PREFIX)
	sb.WriteString(strconv.Itoa(account.id))
	account.name = sb.String()
	mgr.idx += 1
	return nil, account
}