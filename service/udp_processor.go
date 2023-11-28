package service

import (
	"fmt"
	"log"
	"net"
	"spiritDNS/network"
)

type UDPProcessor struct {
	addr *net.UDPAddr
	ctx  *Context
}

func (p *UDPProcessor) Listen() {
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

	p.HandleConnetion(conn)
}

func (p *UDPProcessor) HandleConnetion(conn net.Conn) {
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
		p.ctx = NewContext(buffer[:n])
		p.addr = addr

		// 处理数据
		//go ProcessReq(p)
		ProcessReq(p)
	}
}

func (p *UDPProcessor) FindAnswers(ipList []string) ([]byte, error) {
	data, err := network.TrySendUDPRequest(ipList, 53, p.ctx.RawData)
	if err != nil {
		return nil, fmt.Errorf("network.TrySendUDPRequest err: %s", err.Error())
	}

	// 处理数据
	dnsHeader := parseDNSHeader(data)
	//fmt.Printf("Response: %q\n", data[:n])
	if dnsHeader.Flags.TC {
		// 重新发起TCP
		data, err = network.TrySendTCPRequest(ipList, 53, p.ctx.RawData)
		if err != nil {
			return nil, fmt.Errorf("network.TrySendTCPRequest err: %s", err.Error())
		}
		// 重新读头
		dnsHeader = parseDNSHeader(data)
	}
	if dnsHeader.Flags.AA {
		// 返回 data
		return data, nil
	}

	_, offset := parseDNSQuestion(data, dnsHeader.QuestionCount)

	answerRecords, offset := parseDNSResourceRecord(data, offset, dnsHeader.AnswerCount)
	fmt.Println(answerRecords)

	NSRecords, offset := parseDNSResourceRecord(data, offset, dnsHeader.NSCount)
	fmt.Println(NSRecords)

	additionalRecords, offset := parseDNSResourceRecord(data, offset, dnsHeader.ARCount)
	fmt.Println(additionalRecords)

	return nil, nil
}
