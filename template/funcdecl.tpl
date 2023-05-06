package msg

import  (
    "github.com/springstar/robot/pb"
    "github.com/springstar/robot/core"
    "github.com/golang/protobuf/proto"
    "github.com/jhump/protoreflect/desc"
    "github.com/jhump/protoreflect/dynamic"   
)

<%="var descriptors map[int32]*desc.MessageDescriptor = make(map[int32]*desc.MessageDescriptor)"%>

func AddDescriptor(id int32, desc *desc.MessageDescriptor) {
    descriptors[id] = desc
}

<%= for (n) in names { %>
<%="func Parse"%><%=n%>(<%="id int32, bytes []byte"%>) *<%="pb."%><%=n%> {
    <%= "msg := "%> &<%="pb."%><%=n%>{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}
<% } %>


<%= for (n) in names { %>
<%="func Serialize"%><%=n%>(msgid uint32, <%=params(n)%>) <%="[]byte"%> {
    msg := &pb.<%=n%>{}
<%= for (f) in fields(n) {%>
    <%="msg." %><%=capitalize(f)%> = <%=f%><% } %>
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}
<% } %>