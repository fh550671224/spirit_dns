package main

import (
	"spiritDNS/service"
)

func main() {
	service.ListenUDP()
	service.PackDispatcher.Init()
}
