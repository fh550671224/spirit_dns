package service

import (
	"fmt"
	"net"
	"spiritDNS/dns"
	"spiritDNS/shared"
)

func Resolve(clientQuery *dns.Msg, hostList []string) (*dns.Msg, error) {
	resp := new(dns.Msg)
	resp.MsgHdr = dns.MsgHdr{
		Id:                 clientQuery.Id,
		Response:           true,
		Opcode:             0,
		Authoritative:      false,
		RecursionDesired:   clientQuery.RecursionDesired,
		RecursionAvailable: true,
		Rcode:              0,
	}
	resp.Question = clientQuery.Question

	if len(clientQuery.Question) == 0 {
		resp.MsgHdr.Rcode = dns.RcodeFormatError
		return resp, nil
	}

	question := clientQuery.Question[0]

	if a, ok := answerCache.Get(question); ok {
		resp.Answer = a.answers
		return resp, nil
	}

	queryMsgData, err := clientQuery.Pack()
	if err != nil {
		return nil, fmt.Errorf("Msg.Pack err: %v", err)
	}

	var addrList []*net.UDPAddr
	for _, v := range hostList {
		addrList = append(addrList, &net.UDPAddr{
			IP:   net.ParseIP(v),
			Port: 53,
		})
	}
	pack, err := TrySendUDP(addrList, queryMsgData)
	if err != nil {
		return nil, fmt.Errorf("TrySendUDP err: %v", err)
	}

	for i := 0; i < shared.MaxLookUpTime; i++ {
		ip := pack.Ip
		port := pack.Port
		msg := pack.DnsMsg

		// 如果有截断，重新发起TCP请求
		if msg.MsgHdr.Truncated {
			pack, err = TrySendTCPRequest([]*net.TCPAddr{{
				IP:   net.ParseIP(ip),
				Port: port,
			}}, queryMsgData)
			if err != nil {
				return nil, fmt.Errorf("TrySendTCPRequest err: %v", err)
			}

			msg = pack.DnsMsg
		}

		// 找到答案了，直接返回
		if len(msg.Answer) > 0 {
			var answers []dns.RR
			temp := make([]dns.RR, len(msg.Answer))
			copy(temp, msg.Answer)

			for len(temp) > 0 {
				ans := temp[0]

				answers = append(answers, ans)
				temp = temp[1:]

				if ans.Header().Rrtype == dns.TypeA || ans.Header().Rrtype == dns.TypeAAAA {
					continue
				}

				if ans.Header().Rrtype == dns.TypeCNAME {
					// 需要查询CName记录里的域名解析
					if Cr, ok := ans.(*dns.CNAME); ok {
						m := new(dns.Msg)
						m.SetQuestion(Cr.Target, dns.TypeA)
						res, err := Resolve(m, hostList)
						if err != nil {
							return nil, err
						}

						answers = append(answers, res.Answer...)
					}

				}
			}

			// 存入cache
			go answerCache.Store(question, answers)

			// 返回结果
			resp.Answer = answers
			return resp, nil
		}

		if len(msg.Ns) > 0 {
			// 只有NS记录
			if len(msg.Extra) > 0 {
				// 有Extra记录，直接使用Extra记录的ip递归查询

				var addrs []*net.UDPAddr
				for _, r := range msg.Extra {
					if Ar, ok := r.(*dns.A); ok {
						addrs = append(addrs, &net.UDPAddr{
							IP:   Ar.A,
							Port: 53,
						})
					}

				}
				pack, err = TrySendUDP(addrs, queryMsgData)
				if err != nil {
					return nil, fmt.Errorf("TrySendUDP err: %v", err)
				}
			} else {
				// 需要查询NS记录里的域名解析
				m := new(dns.Msg)
				m.SetQuestion(msg.Ns[0].Header().Name, dns.TypeA)
				res, err := Resolve(m, []string{pack.Ip})
				if err != nil {
					return nil, err
				}

				var addrs []*net.UDPAddr
				for _, r := range res.Answer {
					if Ar, ok := r.(*dns.A); ok {
						addrs = append(addrs, &net.UDPAddr{
							IP:   Ar.A,
							Port: 53,
						})
					}
				}
				pack, err = TrySendUDP(addrs, queryMsgData)
				if err != nil {
					return nil, fmt.Errorf("TrySendUDP err: %v", err)
				}
			}

		}
	}

	return nil, fmt.Errorf("not resolved")
}
