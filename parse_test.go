package main

import (
	"fmt"
	"testing"
	"github.com/stretchr/testify/assert"
	_ "github.com/golang/protobuf/protoc-gen-go/descriptor"
	_ "github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/descriptor"


	"github.com/springstar/robot/pb"
	_ "github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"

)

func TestParse(t *testing.T) {
	addr := &pb.Address{
		State: "china",
		Province: "fujian",
		City:"xiamen",
		Code: 92,
		User: &pb.User{
			Name: "ruida",
		},
		
	   }
	
	d, _ := proto.Marshal(addr)

	var parser protoparse.Parser
	fds, err := parser.ParseFiles("msg/protocol/test.proto")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, fd := range fds {
		fmt.Print(fd.GetName())

		msg := fd.FindMessage("message.Address")
		if (msg != nil) {
			fmt.Println("\tmsg found")
		}

		dmsg := dynamic.NewMessage(msg)
		if (dmsg != nil) {
			fmt.Println(dmsg.GetMessageDescriptor().GetName())
		}

		dmsg.Unmarshal(d)

		var addrMsg pb.Address
		dmsg.ConvertTo(&addrMsg)

		assert.Equal(t, "china", addrMsg.GetState())
		assert.Equal(t, "fujian", addrMsg.GetProvince())
		assert.Equal(t, "xiamen", addrMsg.GetCity())
		assert.Equal(t, int32(92), addrMsg.GetCode())
		assert.Equal(t, "ruida", addrMsg.GetUser().GetName())

		for _, msgDesc := range fd.GetMessageTypes() {
			fmt.Println("\t", msgDesc.GetName())
			for _, fieldDesc := range msgDesc.GetFields() {
				fmt.Println("\t", fieldDesc.GetType().String(), fieldDesc.GetName())
			}
		}
	} 
}

func TestI2dMsg(t *testing.T) {
	// var id2msg map[int32]desc.MessageDescriptor
	// id2msg = make(map[int32]desc.MessageDescriptor)
	// id2msg[101] = pb.Address

}