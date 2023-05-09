package core

type ISubscriber interface {
	HandleMessage(packet *Packet)
}