package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"spiritDNS/client"
	"spiritDNS/handlers"
	"spiritDNS/service"
)

func main() {
	service.InitPacketDispatcher()
	service.InitCache()

	client.InitRedis()
	defer client.RedisClient.CloseRedis()

	client.InitRabbit()
	defer client.RabbitClient.CloseRabbit()

	// Temporarily disable traditional dns
	//go service.ListenUDP()

	go ListenHttp()

	forever := make(chan bool)
	<-forever
}

func ListenHttp() {
	r := gin.Default()

	r.GET("/dns-query", func(c *gin.Context) {
		code, resp, err := handlers.DNSQueryHandler(c)
		if err != nil {
			c.JSON(code, gin.H{
				"msg": err.Error(),
			})
		}

		c.Data(http.StatusOK, "application/dns-message", resp)
	})

	r.Run()
}
