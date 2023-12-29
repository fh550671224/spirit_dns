package service

import (
	"fmt"
	"log"
	"net"
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
			log.Fatal(err)
		}
		msg := shared.DecodeDNSMessage(buffer[:n])

		// TODO 支持常见的类型
		if msg.Questions[0].Type != shared.TYPE_A {
			continue
		}

		if msg.Header.Flags.QR {
			// 查询，分配一个worker来处理
			go Work(&WorkContext{
				ClientAddr:  addr.IP,
				ClientQuery: msg,
			})
		} else {
			// 答案，分发给相应worker
			PackDispatcher.Dispatch(Packet{
				DnsMsg: msg,
				Ip:     addr.IP,
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

func TrySendUDP(ipList []string, port int, data []byte) error {
	for i := 0; i < len(ipList); i++ {
		ip := ipList[i]
		err := sendUDP(data, &net.UDPAddr{
			IP:   net.ParseIP(ipList[i]),
			Port: port,
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
