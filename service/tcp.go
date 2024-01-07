package service

import (
	"fmt"
	"log"
	"net"
	"spiritDNS/dns"
	"spiritDNS/network"
)

//func ListenTCP() {
//	addr := "0.0.0.0:53"
//
//	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	listener, err := net.ListenTCP("udp", tcpAddr)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer listener.Close()
//
//	for {
//		conn, err := listener.Accept()
//	}
//
//	fmt.Printf("Listening on %v\n", addr)
//
//	HandleConnectionTCP(*conn)
//}
//
//func HandleConnectionTCP(conn net.TCPConn) {
//	buffer := make([]byte, 1024)
//	for {
//		n, addr, err := conn.ReadFromUDP(buffer)
//		if err != nil {
//			log.Fatal(err)
//		}
//		msg := shared.DecodeDNSMessage(buffer[:n])
//
//		// TODO 支持常见的类型
//		if msg.Questions[0].QType != shared.TYPE_A {
//			continue
//		}
//
//		if msg.Header.Flags.QR {
//			// 查询，分配一个worker来处理
//			go service.Work(&service.WorkContext{
//				ClientAddr:  addr.IP,
//				ClientPort:  addr.Port,
//				ClientQuery: *msg,
//			})
//		} else {
//			// 答案，分发给相应worker
//			service.MsgDispatcher.Dispatch(*msg)
//		}
//	}
//}
//
//func sendTCP(data []byte, addr *net.TCPAddr) error {
//	tcpConn, err := net.DialTCP("tcp", nil, addr)
//	if err != nil {
//		return fmt.Errorf("net.DialTCP err: %v",err)
//
//	}
//	//defer tcpConn.Close()
//
//	dataLength := len(data)
//	tcpData := make([]byte, 2+dataLength)
//	tcpData[0] = byte(dataLength >> 8)
//	tcpData[1] = byte(dataLength & 0xff)
//	copy(tcpData[2:], data)
//
//	_, err = tcpConn.Write(tcpData)
//	if err != nil {
//		return fmt.Errorf("net.Write err: %v",err)
//	}
//
//	return nil
//}
//
//func TrySendTCP(ipList []string, port int, data []byte) error {
//	for i := 0; i < len(ipList); i++ {
//		ip := ipList[i]
//		err := sendTCP(data, &net.TCPAddr{
//			IP:   net.ParseIP(ipList[i]),
//			Port: port,
//		})
//
//		if err != nil {
//			log.Printf("%s sendTCP err: %v, trying next ip", ip, err)
//		} else {
//			log.Printf("%s sendTCP success", ip)
//			return nil
//		}
//	}
//
//	return fmt.Errorf("no available ip")
//}

func TrySendTCPRequest(addrList []*net.TCPAddr, data []byte) (*Packet, error) {
	for _, addr := range addrList {
		resp, err := network.SendTCP(data, addr)
		if err != nil {
			log.Printf("%s sendTCPRequest err: %v, trying next addr", addr.String(), err)
			continue
		}
		log.Printf("%s sendTCPRequest success", addr.String())

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

	return nil, fmt.Errorf("no available ip")
}
