package service

import (
	"fmt"
	"log"
	"net"
	"spiritDNS/shared"
)

type WorkContext struct {
	ClientAddr  net.IP
	ClientQuery shared.DNSMessage
}

func Work(ctx *WorkContext) {
	answer, err := Resolve(ctx)
	if err != nil {
		log.Printf("Resolve err: %s\n", err.Error())
	}

	// 返回结果给客户端
	encoded, err := shared.EncodeDNSMessage(answer, true)
	if err != nil {
		log.Printf("shared.EncodeDNSMessage err: %s\n", err.Error())
	}

	err = TrySendUDP([]string{ctx.ClientAddr.String()}, 53, encoded)
	if err != nil {
		log.Printf("TrySendUDP err: %s\n", err.Error())
	}
}

func Resolve(ctx *WorkContext) (*shared.DNSMessage, error) {
	encoded, err := shared.EncodeDNSMessage(&ctx.ClientQuery, true)
	if err != nil {
		return nil, fmt.Errorf("shared.EncodeDNSMessage err: %s", err.Error())
	}

	err = TrySendUDP(shared.ROOT_DNS_SERVERS, 53, encoded)
	if err != nil {
		return nil, fmt.Errorf("TrySendUDP err: %s", err.Error())
	}

	ch := make(chan Packet)
	// TODO 支持multi-question
	PackDispatcher.Register(*ctx.ClientQuery.Questions[0], ch)

	for {
		select {
		case pack := <-ch:
			msg := pack.DnsMsg
			ip := pack.Ip

			fmt.Println(pack)
			if msg.Header.Flags.TC {
				// 截断，发起TCP请求
				data, err := TrySendTCPRequest([]string{ip.String()}, 53, encoded)
				if err != nil {
					return nil, fmt.Errorf("TrySendTCPRequest err: %s", err.Error())
				}
				msg = shared.DecodeDNSMessage(data)
			}
		}
	}
}
