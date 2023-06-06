package server

type StatType int32

const (

)

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

func (rs *RunStat) statistic(s Stat) {
	
}




