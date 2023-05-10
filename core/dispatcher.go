package core

import (
	_ "fmt"
)

type IDispatcher interface {
	Register(uint32, ISubscriber)
	Dispatch(*Packet)
}

type MessageDispatcher struct {
	subscribers map[uint32][]ISubscriber
}

func NewMsgDispatcher() *MessageDispatcher {
	return &MessageDispatcher{
		subscribers: make(map[uint32][]ISubscriber),
	}
}

func (d *MessageDispatcher) Register(msgid uint32, subscriber ISubscriber) {
	subs, ok := d.subscribers[msgid]
	if ok {
		subs = append(subs, subscriber)
	} else {
		var ls []ISubscriber
		ls = append(ls, subscriber)
		d.subscribers[msgid] = ls
	}
}

func (d *MessageDispatcher) Dispatch(packet *Packet) {
	msgid := packet.Type
	subs, ok := d.subscribers[msgid]
	if !ok {
		return
	}

	for _, subscriber := range subs {
		subscriber.HandleMessage(packet)
	}
}