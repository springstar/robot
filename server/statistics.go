package server

type ServerStat struct {
	sends  		int32 `json:"sends"`
	roles  		int32 `json:"creates"`
	logins 		int32 `json:"logins"`
	spackets 	int32 `json:"sendpackets"`
	rpackets	int32 `json:"recvpackets"`
	sbytes		int32 `json:"sendbytes"`
	rbytes		int32 `json:"recvbytes"`
	msgsends   map[int32]int32	`json:"msgsends"`
	msgrecvs   map[int32]int32	`json:"msgrecvs"`
}

func newServerStat() *ServerStat {
	return &ServerStat{
		msgsends: make(map[int32]int32),
		msgrecvs: make(map[int32]int32),
	}
}

