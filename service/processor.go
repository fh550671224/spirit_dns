package service

import (
	"net"
)

type Context struct {
}

type Processor interface {
	Listen()
	HandleConnection(conn net.Conn)
	FindAnswers(ipList []string) ([]byte, error)
}

func StarListening() {
	udpp := UdpProcessor{}
	udpp.Listen()
}

func ProcessReq(p Processor) {
	// TODO 写
}
