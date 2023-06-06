package server

type StatType int32

const (
	STAT_SEND_ROBOTS = iota
	STAT_CREATE_ROLES
	STAT_ENTER_STAGE
	STAT_SEND_PACKETS
	STAT_RECV_PACKETS
	STAT_SEND_BYTES
	STAT_RECV_BYTES

)

type MsgStat struct {
	msgTyp int32
	bytes int32
}

type Stat struct {
	typ StatType
	v interface{}
}

type RunStat struct {
	sends  		int32 `json:"sends"`
	roles  		int32 `json:"creates"`
	logins 		int32 `json:"logins"`
	spackets 	int32 `json:"sendpackets"`
	rpackets	int32 `json:"recvpackets"`
	sbytes		int32 `json:"sendbytes"`
	rbytes		int32 `json:"recvbytes"`
	msgsends   map[int32]int32	`json:"msgsends"`
	msgrecvs   map[int32]int32	`json:"msgrecvs"`

	ch chan Stat

}

func newRunStat() *RunStat {
	return &RunStat{
		msgsends: make(map[int32]int32),
		msgrecvs: make(map[int32]int32),
		ch: make(chan Stat),
	}
}

func (rs *RunStat) queueMsgStat(stype StatType, mtype int32, size int32) {
	var s Stat
	s.typ = stype
	var ms MsgStat
	ms.msgTyp = mtype
	ms.bytes = size
	s.v = ms
	rs.ch <- s
}

func (rs *RunStat) statistic(s Stat) {
	switch s.typ {
	case STAT_SEND_PACKETS:
		rs.statSendPackets(s)
	case STAT_RECV_PACKETS:	
		rs.statRecvPackets(s)
	default:

	}
}

func (rs *RunStat) statSendPackets(s Stat) {

}

func (rs *RunStat) statRecvPackets(s Stat) {

}



