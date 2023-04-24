package core

type NetConnection interface {
	Connect(addr string) error
}



