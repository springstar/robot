package main

import  (
    "github.com/springstar/robot/pb"
    "github.com/jhump/protoreflect/desc"
    "github.com/jhump/protoreflect/dynamic"   
)

var descriptors map[int32]*desc.MessageDescriptor = make(map[int32]*desc.MessageDescriptor)

func AddDescriptor(id int32, desc *desc.MessageDescriptor) {
    descriptors[id] = desc
}


func parseCSLogin(id int32, bytes []byte) *pb.CSLogin {
    msg :=  &pb.CSLogin{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func parseCSTest(id int32, bytes []byte) *pb.CSTest {
    msg :=  &pb.CSTest{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}




func serializeCSLogin(account string, password string, token string, serverId int32, version int32) *pb.CSLogin {
    msg := &pb.CSLogin{}

    msg.Account = account
    msg.Password = password
    msg.Token = token
    msg.ServerId = serverId
    msg.Version = version
    return msg
}

func serializeCSTest(code string) *pb.CSTest {
    msg := &pb.CSTest{}

    msg.Code = code
    return msg
}
