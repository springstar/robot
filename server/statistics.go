package server

import (
	_ "fmt"
	"time"
)

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
	Sends  		int32 `json:"sends"`
	Roles  		int32 `json:"creates"`
	Logins 		int32 `json:"logins"`
	Spackets 	int32 `json:"sendpackets"`
	Rpackets	int32 `json:"recvpackets"`
	Sbytes		int32 `json:"sendbytes"`
	Rbytes		int32 `json:"recvbytes"`
	Msgsends   map[int32]int32	`json:"msgsends"`
	Msgrecvs   map[int32]int32	`json:"msgrecvs"`

	ch chan Stat
	ticker *time.Ticker


}

func newRunStat() *RunStat {
	return &RunStat{
		Msgsends: make(map[int32]int32),
		Msgrecvs: make(map[int32]int32),
		ch: make(chan Stat, 8192),
		ticker : time.NewTicker(SERVER_PULSE),

	}
}

func queueStat(stype StatType, count int32) {
	var s Stat
	s.typ = stype
	s.v = count
	serv.ch <- s

}

func queueMsgStat(stype StatType, mtype int32, size int32) {
	var s Stat
	s.typ = stype
	var ms MsgStat
	ms.msgTyp = mtype
	ms.bytes = size
	s.v = ms
	serv.ch <- s
}

func (rs *RunStat) Start() {
	for {
		select {
		case <- rs.ticker.C:
			rs.pulse()
			rs.ticker.Reset(SERVER_PULSE)	
		}	
	}
}

func (rs *RunStat) pulse() {
	for stat := range rs.ch {
		rs.statistic(stat)
	}
}


func (rs *RunStat) statistic(s Stat) {
	switch s.typ {
	case STAT_SEND_PACKETS:
		rs.statSendPackets(s)
	case STAT_RECV_PACKETS:	
		rs.statRecvPackets(s)
	case STAT_SEND_ROBOTS:
		rs.Sends += s.v.(int32)
	case STAT_CREATE_ROLES:
		rs.Roles += s.v.(int32)	
	case STAT_ENTER_STAGE:	
		rs.Logins += s.v.(int32)
	default:
		break
	}
}

func (rs *RunStat) statSendPackets(s Stat) {
	rs.Spackets += 1
	ms := s.v.(MsgStat)
	rs.Sbytes += ms.bytes
	rs.Msgsends[ms.msgTyp] = rs.Msgsends[ms.msgTyp] + 1
}

func (rs *RunStat) statRecvPackets(s Stat) {
	rs.Rpackets += 1
	ms := s.v.(MsgStat)
	rs.Rbytes += ms.bytes

	rs.Msgrecvs[ms.msgTyp] = rs.Msgrecvs[ms.msgTyp] + 1
}



