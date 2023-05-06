package core

import (
	_ "fmt"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	msg := []byte{10, 5, 14, 111, 98, 111, 116, 18, 6, 49, 50, 51, 52, 53, 54, 32, 233, 7, 40, 1}
	packet := NewPacket(111, msg)
	b := packet.Serialize()
	p := Parse(b)
	assert.Equal(t, uint32(111), p.Type)
	assert.Equal(t, uint32(len(msg) + 8), p.Length)

}