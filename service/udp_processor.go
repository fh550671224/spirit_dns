package service

import (
	"fmt"
	"log"
	"net"
	"reflect"
	"spiritDNS/network"
	"spiritDNS/shared"
)

type UdpProcessor struct {
	addr *net.UDPAddr
	msg  *shared.DNSMessage
}

func (p *UdpProcessor) Listen() {
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

	p.HandleConnection(conn)
}

func (p *UdpProcessor) HandleConnection(conn net.Conn) {
	udpConn, ok := conn.(*net.UDPConn)
	if !ok {
		log.Fatal("Failed to convert to *net.UDPConn")
	}
	buffer := make([]byte, 1024)
	for {
		n, addr, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal(err)
		}

		p.addr = addr
		p.msg = shared.DecodeDNSMessage(buffer[:n])

		if p.msg.Questions[0].Type != shared.TYPE_A {
			continue
		}

		// 处理数据 TODO 改成goroutine
		//go ProcessReq(p)
		p.FindAnswers(shared.ROOT_DNS_SERVERS)
	}
}

func (p *UdpProcessor) FindAnswers(ipList []string) ([]byte, error) {
	encoded, err := shared.EncodeDNSMessage(p.msg, true)
	if err != nil {
		return nil, fmt.Errorf("shared.EncodeDNSMessage err: %s", err.Error())
	}

	data, err := network.TrySendUDPRequest(ipList, 53, encoded)
	if err != nil {
		return nil, fmt.Errorf("network.TrySendUDPRequest err: %s", err.Error())
	}

	// TODO 处理数据
	msg := shared.DecodeDNSMessage(data)

	// TODO for test, delete
	{
		encoded, err = shared.EncodeDNSMessage(msg, true)
		if err != nil {
			log.Fatal(err)
		}
		//if !bytes.Equal(data, encoded) {
		//	log.Fatal("encoded not equal to decoded")
		//}

		msgg := shared.DecodeDNSMessage(encoded)

		a := reflect.DeepEqual(msg, msgg)

		fmt.Println(a)
	}

	return nil, nil
}
