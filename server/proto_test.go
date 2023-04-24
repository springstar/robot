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
	message := &pb.CSTest{
		Code : "10054",
	}

	packet, _ := proto.Marshal(message)

	msg := &pb.CSTest{}
	proto.Unmarshal(packet, msg)
	assert.Equal(t, "10054", msg.Code)
	assert.Equal(t, "10054", msg.GetCode())

	_, md := descriptor.ForMessage(msg)
	options := md.GetOptions()
	a, _ := proto.GetExtension(options, pb.E_Msgid)
   assert.Equal(t, int32(101), *a.(*int32))

   addr := &pb.Adress{
	State: "china",
	Province: "fujian",
	City:"xiamen",
	Code: 92,
	User: &pb.User{
		Name: "ruida",
	},
	
   }

   d, _ := proto.Marshal(addr)
   ad := &pb.Adress{}
   proto.Unmarshal(d, ad)
   assert.Equal(t, "china", ad.GetState())
   assert.Equal(t, "fujian", ad.GetProvince())
   assert.Equal(t, "xiamen", ad.GetCity())
   assert.Equal(t, int32(92), ad.GetCode())
   assert.Equal(t, "ruida", ad.GetUser().GetName())

	
	
}