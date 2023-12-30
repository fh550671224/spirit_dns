package service

import (
	"fmt"
	"github.com/miekg/dns"
	"log"
	"net"
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

		if msg.MsgHdr.Response {
			// 答案，分发给相应worker
			PackDispatcher.Dispatch(Packet{
				DnsMsg: *msg,
				Ip:     addr.IP,
			})
		} else {
			// 查询，分配一个worker来处理

			// TODO 添加其他Opcode
			if msg.Opcode != 0 {
				continue
			}

			go Work(&WorkContext{
				ClientAddr:  addr.IP,
				ClientQuery: *msg,
			})
		}
	}
}

func sendUDP(data []byte, addr *net.UDPAddr) error {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return fmt.Errorf("net.DialUDP err: %s", err.Error())

	}
	defer conn.Close()

	// send request
	_, err = conn.Write(data)
	if err != nil {
		return fmt.Errorf("net.Write err: %s", err.Error())
	}

	return nil
}

func TrySendUDP(ipList []string, data []byte) error {
	for i := 0; i < len(ipList); i++ {
		ip := ipList[i]
		err := sendUDP(data, &net.UDPAddr{
			IP:   net.ParseIP(ipList[i]),
			Port: 53,
		})

		if err != nil {
			log.Printf("%s sendUDPRequest err: %v, trying next ip", ip, err)
		} else {
			log.Printf("%s sendUDPRequest success", ip)
			return nil
		}
	}
	return fmt.Errorf("no available ip")
}
