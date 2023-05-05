package core

import (
	"encoding/binary"
	"io"
	"bytes"
)

// Send buffer size determines how many bytes we send in a single TCP write call.
// This can be anything from 1 to 65495.
// A good default value for this can be read from: /proc/sys/net/ipv4/tcp_wmem
const sendBufferSize = 16384

// Packet represents a single network message.
// It has a byte code indicating the type of the message
// and a data payload in the form of a byte slice.
type Packet struct {
	Length int32
	Msgid  int32
	Data   []byte
}

// New creates a new packet.
// It expects a byteCode for the type of message and
// a data parameter in the form of a byte slice.
func New(msgid int32, data []byte) *Packet {
	return &Packet{
		Msgid:   msgid,
		Length: int32(len(data)) + 8,
		Data:   data,
	}
}

// Write writes the packet to the IO device.
func (packet *Packet) Write(writer io.Writer) error {
	err := binary.Write(writer, binary.BigEndian, packet.Length)

	if err != nil {
		return err
	}

	err = binary.Write(writer, binary.BigEndian, packet.Msgid)

	if err != nil {
		return err
	}



	n := 0
	bytesWritten := 0
	writeUntil := 0

	for bytesWritten < len(packet.Data) {
		writeUntil = bytesWritten + sendBufferSize

		if writeUntil > len(packet.Data) {
			writeUntil = len(packet.Data)
		}

		n, err = writer.Write(packet.Data[bytesWritten:writeUntil])

		if err != nil {
			return err
		}

		bytesWritten += n
	}

	return err
}

// Bytes returns the raw byte slice serialization of the packet.
func (packet *Packet) Bytes() []byte {
	result := []byte{}
	result = append(result, Int32ToBytes(packet.Msgid)...)
	result = append(result, Int32ToBytes(packet.Length)...)
	result = append(result, packet.Data...)
	return result
}

func Int32ToBytes(i int32) []byte {
	buffer := bytes.Buffer{}
	_ = binary.Write(&buffer, binary.BigEndian, i)
	return buffer.Bytes()
}