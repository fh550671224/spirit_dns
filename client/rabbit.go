package client

import (
	dns "github.com/fh550671224/spirit_dns_public"
	"log"
)

var RabbitClient dns.RabbitClient

type LogMsg struct {
	QuerySource int
	Request     interface{}
	Addr        string
	Data        dns.Msg
}

func InitRabbit() {
	err := RabbitClient.Init("guest", "guest", "localhost:5672")
	if err != nil {
		log.Printf("RabbitMQ Init err:%v", err)
	}
}
