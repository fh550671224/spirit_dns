package main

import (
	"spiritDNS/client"
	"spiritDNS/service"
)

func main() {
	service.InitPacketDispatcher()
	service.InitCache()

	client.InitRedis()
	defer client.RedisClient.CloseRedis()

	service.ListenUDP()
}
