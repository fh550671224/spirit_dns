package client

import (
	dns "github.com/fh550671224/spirit_dns_public"
	"log"
)

var RabbitClient dns.RabbitClient

func InitRabbit() {
	err := RabbitClient.Init("guest", "guest", "localhost:5672")
	if err != nil {
		log.Printf("RabbitMQ Init err:%v", err)
	}
}
