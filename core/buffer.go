package core

import (
	_ "fmt"
	"bytes"
	_ "errors"
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
	if b.buf.Len() < 4 {
		return nil
	}

	var packets []*Packet

	for {
		data := b.buf.Bytes()
		len, err := uint32FromBytes(data[0:4])
		if b.buf.Len() < int(len) || len == 0 {
			return packets
		}

		slice := make([]byte, len)
		b.buf.Read(slice)

		msgid, err := uint32FromBytes(slice[4:8])
		if (err != nil) {
			return packets
		}
	
		packet := NewPacket(msgid, slice[8:])
		packets = append(packets, packet)
	}

	return packets
}