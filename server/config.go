package server

type ServerConfig struct {
	ServerId int	`yaml:"server"`
	MaxNum int	`yaml:"maxnum"`
	Url string `yaml:"url"`
}

func  getStartAccountId(serverId int) int {
	return serverId * ROBOT_SECTION
}