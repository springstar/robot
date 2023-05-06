package core

import (
	
)

type NetConnection interface {
	Connect(addr string) error
	Write(p []byte) (n int, err error)
	Read() ([]byte, error)
}



