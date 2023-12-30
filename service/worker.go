package service

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
	"spiritDNS/shared"
)

type WorkContext struct {
	ClientAddr  net.IP
	ClientQuery dns.Msg
}

func Work(ctx *WorkContext) {
	answer, err := Resolve(ctx.ClientQuery)
	if err != nil {
		log.Printf("Resolve err: %s\n", err.Error())
	}

	// 返回结果给客户端
	answerData, err := answer.Pack()
	if err != nil {
		log.Printf("shared.EncodeDNSMessage err: %s\n", err.Error())
	}

	err = TrySendUDP([]string{ctx.ClientAddr.String()}, answerData)
	if err != nil {
		log.Printf("TrySendUDP err: %s\n", err.Error())
	}
}

func Resolve(clientQuery dns.Msg) (*dns.Msg, error) {
	queryMsgData, err := clientQuery.Pack()
	if err != nil {
		return nil, fmt.Errorf("Msg.Pack err: %s", err.Error())
	}

	err = TrySendUDP(shared.ROOT_DNS_SERVERS, queryMsgData)
	if err != nil {
		return nil, fmt.Errorf("TrySendUDP err: %s", err.Error())
	}

	ch := make(chan Packet)
	PackDispatcher.Register(clientQuery.Question[0].Name, ch)
	defer PackDispatcher.UnRegister(clientQuery.Question[0].Name)

	for {
		select {
		case pack := <-ch:
			msg := pack.DnsMsg
			ip := pack.Ip

			if msg.MsgHdr.Truncated {
				// 截断，发起TCP请求
				data, err := TrySendTCPRequest([]string{ip.String()}, 53, queryMsgData)
				if err != nil {
					return nil, fmt.Errorf("TrySendTCPRequest err: %s", err.Error())
				}
				err = msg.Unpack(data)
				if err != nil {
					return nil, fmt.Errorf("Msg.Unpack err: %s", err.Error())
				}
			}

			if len(msg.Answer) > 0 {
				// 找到答案了，直接返回
				return &msg, nil
			}

			if len(msg.Ns) > 0 {
				// 只有NS记录
				if len(msg.Extra) > 0 {
					// 有Extra记录，直接使用Extra记录的ip递归查询

					var addrs []string
					for _, r := range msg.Extra {
						addrs = append(addrs, r.Header().Name)
					}
					err = TrySendUDP(addrs, queryMsgData)
					if err != nil {
						return nil, fmt.Errorf("TrySendUDP err: %s", err.Error())
					}

					// 发送完请求继续阻塞等待

					//ch = make(chan Packet)
					//PackDispatcher.Register(ctx.ClientQuery.Question[0], ch)
				} else {
					// 需要查询NS记录里的域名解析
					m := new(dns.Msg)
					m.SetQuestion(msg.Ns[0].Header().Name, dns.TypeA)
					answer, err := Resolve(*m)
					if err != nil {
						return nil, err
					}

					var addrs []string
					for _, r := range answer.Answer {
						addrs = append(addrs, r.Header().Name)
					}
					err = TrySendUDP(addrs, queryMsgData)
					if err != nil {
						return nil, fmt.Errorf("TrySendUDP err: %s", err.Error())
					}
				}

			}
		}
	}
}
