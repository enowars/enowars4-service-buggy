package Mailing

type mailConf struct {
	Host string
	Port int
}

func SetMailConfig(host string, port int) *mailConf {
	m := &mailConf{
		Host: host,
		Port: port,
	}
	return m
}
