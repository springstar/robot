package core

import (
	_ "fmt"
	"encoding/binary"
	_ "io"
	"bytes"
)


const sendBufferSize = 16384


type Packet struct {
	Length uint32
	Type  uint32
	Data   []byte
}

func NewPacket(msgid uint32, data []byte) *Packet {
	return &Packet{
		Length: uint32(len(data)) + 8,
		Type:   msgid,
		Data:   data,
	}
}

func (packet *Packet) Serialize() []byte {
	result := []byte{}
	result = append(result, uint32ToBytes(packet.Length)...)
	result = append(result, uint32ToBytes(packet.Type)...)
	result = append(result, packet.Data...)
	return result
}

func uint32ToBytes(i uint32) []byte {
	buffer := bytes.Buffer{}
	_ = binary.Write(&buffer, binary.BigEndian, i)
	return buffer.Bytes()
}

func uint32FromBytes(b []byte) (uint32, error) {
	buffer := bytes.NewReader(b)
	var result uint32
	err := binary.Read(buffer, binary.BigEndian, &result)
	return result, err
}

