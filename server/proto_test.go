package server


import (
	"testing"
	"github.com/springstar/robot/pb"
	"github.com/stretchr/testify/assert"
	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
)

func TestMsgId(t *testing.T) {
	var msg *pb.CSTest
	_, md := descriptor.ForMessage(msg)
	options := md.GetOptions()
	a, _ := proto.GetExtension(options, pb.E_Msgid)
   assert.Equal(t, int32(101), *a.(*int32))
}

func TestMarshal(t *testing.T) {
	
}