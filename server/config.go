package server

type ServerConfig struct {
	ServerId int	`yaml:"server"`
	MaxNum int	`yaml:"maxnum"`
}

func  getStartAccountId(serverId int) int {
	return serverId * ROBOT_SECTION
}