package msg

import  (
    "github.com/springstar/robot/pb"
    "github.com/springstar/robot/core"
    "github.com/golang/protobuf/proto"
    "github.com/jhump/protoreflect/desc"
    "github.com/jhump/protoreflect/dynamic"   
)

var descriptors map[int32]*desc.MessageDescriptor = make(map[int32]*desc.MessageDescriptor)

func AddDescriptor(id int32, desc *desc.MessageDescriptor) {
    descriptors[id] = desc
}


func ParseCSQueryCharacters(id int32, bytes []byte) *pb.CSQueryCharacters {
    msg :=  &pb.CSQueryCharacters{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseCSCharacterCreate(id int32, bytes []byte) *pb.CSCharacterCreate {
    msg :=  &pb.CSCharacterCreate{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseSCCharacterLoginResult(id int32, bytes []byte) *pb.SCCharacterLoginResult {
    msg :=  &pb.SCCharacterLoginResult{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseSCInitData(id int32, bytes []byte) *pb.SCInitData {
    msg :=  &pb.SCInitData{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseCSLogin(id int32, bytes []byte) *pb.CSLogin {
    msg :=  &pb.CSLogin{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseSCQueryCharactersResult(id int32, bytes []byte) *pb.SCQueryCharactersResult {
    msg :=  &pb.SCQueryCharactersResult{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseSCCharacterCreateResult(id int32, bytes []byte) *pb.SCCharacterCreateResult {
    msg :=  &pb.SCCharacterCreateResult{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseCSCharacterLogin(id int32, bytes []byte) *pb.CSCharacterLogin {
    msg :=  &pb.CSCharacterLogin{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseCSStageEnter(id int32, bytes []byte) *pb.CSStageEnter {
    msg :=  &pb.CSStageEnter{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseSCStageEnterResult(id int32, bytes []byte) *pb.SCStageEnterResult {
    msg :=  &pb.SCStageEnterResult{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseCSTest(id int32, bytes []byte) *pb.CSTest {
    msg :=  &pb.CSTest{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseCSPing(id int32, bytes []byte) *pb.CSPing {
    msg :=  &pb.CSPing{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}

func ParseSCLoginResult(id int32, bytes []byte) *pb.SCLoginResult {
    msg :=  &pb.SCLoginResult{}
    md := descriptors[id]
    dmsg := dynamic.NewMessage(md)
    dmsg.Unmarshal(bytes)
    dmsg.ConvertTo(msg)
    return msg
}




func SerializeCSQueryCharacters(msgid uint32, serverId int32) []byte {
    msg := &pb.CSQueryCharacters{}

    msg.ServerId = serverId
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeCSCharacterCreate(msgid uint32, name string, roleSn int32, ptest bool, fashionSn []int32, serverId int32) []byte {
    msg := &pb.CSCharacterCreate{}

    msg.Name = name
    msg.RoleSn = roleSn
    msg.Ptest = ptest
    msg.FashionSn = fashionSn
    msg.ServerId = serverId
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeSCCharacterLoginResult(msgid uint32, resultCode int32) []byte {
    msg := &pb.SCCharacterLoginResult{}

    msg.ResultCode = resultCode
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeSCInitData(msgid uint32, human *pb.DHuman, stage *pb.DInitDataStage, skill *pb.DSkill, skillB *pb.DSkill, sourType pb.ESoulType, dBag *pb.DBag, altar *pb.DBag, tasks *pb.DTask, excuteTaskCode int32, mails *pb.DMail, treasures []*pb.DTreasure, blackGoods []*pb.DBlackGoods, signIn *pb.DSignIn, setting *pb.DSetting, campInfo *pb.DCampInfo, questInfo []*pb.DQuest, shortcuts int64, scpotion []*pb.ShortcutPotion, skillSlot []int32, eCollegeQuestionType pb.ECollegeQuestionType, collegeQuestion []*pb.DCollegeQuestion, startServerTime int64) []byte {
    msg := &pb.SCInitData{}

    msg.Human = human
    msg.Stage = stage
    msg.Skill = skill
    msg.SkillB = skillB
    msg.SourType = sourType
    msg.DBag = dBag
    msg.Altar = altar
    msg.Tasks = tasks
    msg.ExcuteTaskCode = excuteTaskCode
    msg.Mails = mails
    msg.Treasures = treasures
    msg.BlackGoods = blackGoods
    msg.SignIn = signIn
    msg.Setting = setting
    msg.CampInfo = campInfo
    msg.QuestInfo = questInfo
    msg.Shortcuts = shortcuts
    msg.Scpotion = scpotion
    msg.SkillSlot = skillSlot
    msg.ECollegeQuestionType = eCollegeQuestionType
    msg.CollegeQuestion = collegeQuestion
    msg.StartServerTime = startServerTime
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeCSLogin(msgid uint32, account string, password string, token string, serverId int32, version int32) []byte {
    msg := &pb.CSLogin{}

    msg.Account = account
    msg.Password = password
    msg.Token = token
    msg.ServerId = serverId
    msg.Version = version
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeSCQueryCharactersResult(msgid uint32, characters []*pb.DCharacter) []byte {
    msg := &pb.SCQueryCharactersResult{}

    msg.Characters = characters
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeSCCharacterCreateResult(msgid uint32, resultCode int32, humanId int64, resultReason string, fashionSn int32) []byte {
    msg := &pb.SCCharacterCreateResult{}

    msg.ResultCode = resultCode
    msg.HumanId = humanId
    msg.ResultReason = resultReason
    msg.FashionSn = fashionSn
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeCSCharacterLogin(msgid uint32, humanId int64) []byte {
    msg := &pb.CSCharacterLogin{}

    msg.HumanId = humanId
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeCSStageEnter(msgid uint32, ) []byte {
    msg := &pb.CSStageEnter{}

    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeSCStageEnterResult(msgid uint32, obj []*pb.DStageObject, isShowBossNPC bool, isBridge bool) []byte {
    msg := &pb.SCStageEnterResult{}

    msg.Obj = obj
    msg.IsShowBossNPC = isShowBossNPC
    msg.IsBridge = isBridge
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeCSTest(msgid uint32, code string) []byte {
    msg := &pb.CSTest{}

    msg.Code = code
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeCSPing(msgid uint32, ) []byte {
    msg := &pb.CSPing{}

    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}

func SerializeSCLoginResult(msgid uint32, resultCode int32, resultReason string, keyActivate bool, isServerFull bool, showGiftCode int32, showService int32, showGm int32) []byte {
    msg := &pb.SCLoginResult{}

    msg.ResultCode = resultCode
    msg.ResultReason = resultReason
    msg.KeyActivate = keyActivate
    msg.IsServerFull = isServerFull
    msg.ShowGiftCode = showGiftCode
    msg.ShowService = showService
    msg.ShowGm = showGm
    buf, _ := proto.Marshal(msg)
    packet := core.NewPacket(msgid, buf)
    return packet.Serialize()
}
