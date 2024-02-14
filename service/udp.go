package service

import (
	"fmt"
	"github.com/fh550671224/spirit_dns_public"
	"log"
	"net"
	"spiritDNS/network"
	"spiritDNS/shared"
)

func ListenUDP() {
	addr := "0.0.0.0:53"

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Printf("Listening on %v\n", addr)

	HandleConnectionUDP(*conn)
}

func HandleConnectionUDP(conn net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err, addr)
		}

		msg := new(dns.Msg)
		err = msg.Unpack(buffer[:n])
		if err != nil {
			log.Fatal(err)
		}

		// TODO 支持所有opcode
		if msg.Opcode != 0 {
			continue
		}

		//  交给goroutine异步解析
		go func() {
			answer, err := Resolve(msg, shared.ROOT_DNS_SERVERS)
			if err != nil {
				log.Printf("Resolve err: %v\n", err)
				return
			}

			// 返回结果给客户端
			answerData, err := answer.Pack()
			if err != nil {
				log.Printf("dns.Pack err: %v\n", err)
			}

			dispatcher.Dispatch(addr, answerData)
		}()

		// 开启goroutine阻塞等待结果
		go func() {
			ch := make(chan []byte)
			dispatcher.Register(addr, ch)
			defer dispatcher.UnRegister(addr)

			select {
			case data := <-ch:
				_, err = conn.WriteToUDP(data, addr)
				if err != nil {
					log.Fatalf("net.WriteToUDP err: %v", err)
				}
				return
			}
		}()

	}
}

func TrySendUDP(addrList []*net.UDPAddr, data []byte) (*Packet, error) {
	for _, addr := range addrList {
		resp, err := network.SendUDP(data, addr)
		if err != nil {
			log.Printf("%s sendUDPRequest err: %v, trying next addr", addr.String(), err)
			continue
		}

		log.Printf("%s sendUDPRequest success", addr.String())

		msg := new(dns.Msg)
		err = msg.Unpack(resp)
		if err != nil {
			return nil, fmt.Errorf("dns.Unpack err: %v", err)
		}

		if msg.IsInvalid() {
			log.Printf("%s has no Answer or Ns, trying next addr", addr.String())
			continue
		}

		return &Packet{
			DnsMsg: msg,
			Ip:     addr.IP.String(),
			Port:   addr.Port,
		}, nil
	}

	return nil, fmt.Errorf("no available addr")
}
