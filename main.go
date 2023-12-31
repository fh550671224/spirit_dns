package main

import (
	"spiritDNS/service"
)

func main() {
	service.InitPacketDispatcher()
	service.InitCache()

	service.ListenUDP()

}
