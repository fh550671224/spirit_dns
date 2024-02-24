package handlers

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	dns "github.com/fh550671224/spirit_dns_public"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"spiritDNS/client"
	"spiritDNS/service"
	"spiritDNS/shared"
)

func DNSQueryHandler(c *gin.Context) (int, []byte, error) {
	encodedData := c.Query("dns")
	if encodedData == "" {
		return http.StatusBadRequest, nil, fmt.Errorf("query dns not found")
	}

	data, err := base64.RawURLEncoding.DecodeString(encodedData)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("DecodeString err: %v", err)
	}

	msg := new(dns.Msg)
	err = msg.Unpack(data)
	if err != nil {
		return http.StatusBadRequest, nil, fmt.Errorf("msg.Unpack err:%v", err)
	}

	// TODO 支持所有opcode
	if msg.Opcode != 0 {
		return http.StatusNotImplemented, nil, fmt.Errorf("not implemented")
	}

	answer, err := service.Resolve(msg, shared.ROOT_DNS_SERVERS)
	if err != nil {
		return http.StatusInternalServerError, nil, fmt.Errorf("Resolve err: %v", err)
	}

	// 返回结果给客户端
	answerData, err := answer.Pack()
	if err != nil {
		log.Printf("dns.Pack err: %v\n", err)
	}

	// log
	go func() {
		var dnsMsg dns.Msg
		err = dnsMsg.Unpack(data)
		if err != nil {
			log.Printf("dnsMsg.Unpack err:%v", err)
			return
		}

		// 记录日志
		logMsg := client.LogMsg{
			QuerySource: shared.DNS_QUERY_SOURCE_HTTP,
			Request:     c.Request,
			Data:        dnsMsg,
		}
		logBytes, err := json.Marshal(&logMsg)
		if err != nil {
			log.Printf("Marshal logMsg err: %v", err)
			return
		}
		err = client.RabbitClient.Write(dns.SpiritDNSLog, logBytes)
		if err != nil {
			log.Printf("Write logMsg to mq err: %v", err)
		}
	}()

	return http.StatusOK, answerData, nil
}
