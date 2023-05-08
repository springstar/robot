package core

import (
	"fmt"
	"bytes"
)

type PacketBuffer struct {
	buf bytes.Buffer
}

func NewBuffer() *PacketBuffer {
	return &PacketBuffer{
		
	}
}

func (b *PacketBuffer) Write(bytes []byte) error {
	_, err := b.buf.Write(bytes)
	return err
}

func (b *PacketBuffer) Read() []*Packet {
	if b.buf.Len() < 8 {
		return nil
	}

	var packets []*Packet

	for {
		data := b.buf.Bytes()
		fmt.Println(data)
		len, _ := uint32FromBytes(data[0:4])
		if b.buf.Len() < int(len) || len == 0 {
			break
		}
		
		slice := make([]byte, len)
		b.buf.Read(slice)
		packet := Parse(slice)
		packets = append(packets, packet)
	}

	return packets
}