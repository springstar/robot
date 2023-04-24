package core


type Packet struct {
	len 		int32
	msgid 		int32
	payload 	[]byte
}