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

	go service.ListenUDP()

	go ListenHttp()

	forever := make(chan bool)
	<-forever
}

func ListenHttp() {
	// 创建一个Gin路由器
	r := gin.Default()

	// 注册一个GET路由
	// 当访问根URL "/"时，返回"Hello, World!"字符串
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, World!")
	})

	r.GET("/dns-query", func(c *gin.Context) {
		code, resp, err := handlers.DNSQueryHandler(c)
		if err != nil {
			c.JSON(code, gin.H{
				"msg": err.Error(),
			})
		}

		c.Data(http.StatusOK, "application/dns-message", resp)
	})

	// 运行Gin服务器，默认监听在8080端口
	r.Run(":12227") // 监听并在 0.0.0.0:8080 上启动服务
}
