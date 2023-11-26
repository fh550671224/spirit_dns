package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

type DNSQuestion struct {
	Name   string
	QType  uint16
	QClass uint16
}

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	Qdcount uint16
	Ancount uint16
	Nscount uint16
	Arcount uint16
}

func main() {
	addr := "0.0.0.0:53"

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("Listening on %v\n", addr)

	for {
		handleClient(conn)
	}
}

func handleClient(conn *net.UDPConn) {
	buffer := make([]byte, 1024)

	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	// 处理数据
	dnsHeader := parseDNSRequest(buffer[:n])
	//fmt.Printf("Query: %q\n", buffer[:n])

	fmt.Printf("%x %x %x %x %x %x\n", dnsHeader.ID, dnsHeader.Flags,
		dnsHeader.Qdcount, dnsHeader.Ancount, dnsHeader.Nscount, dnsHeader.Arcount)

	questions := parseDNSQuestion(buffer, dnsHeader.Qdcount)

	for _, q := range questions {
		fmt.Println(q.Name)
	}

	//fmt.Printf("length of last name:%x and is: %c %c %c\n", buffer[12], buffer[13], buffer[14], buffer[15])
}

func parseDNSRequest(data []byte) *DNSHeader {
	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		Flags:   binary.BigEndian.Uint16(data[2:4]),
		Qdcount: binary.BigEndian.Uint16(data[4:6]),
		Ancount: binary.BigEndian.Uint16(data[6:8]),
		Nscount: binary.BigEndian.Uint16(data[8:10]),
		Arcount: binary.BigEndian.Uint16(data[10:12]),
	}
}

func parseDNSQuestion(data []byte, Qdcount uint16) []*DNSQuestion {
	offset := 12
	var questions []*DNSQuestion
	for Qdcount > 0 {
		q := DNSQuestion{}

		var names []string
		for data[offset] != 0 {
			length := int(data[offset])
			offset++
			name := data[offset : offset+length]
			names = append(names, string(name))
			offset += length
		}
		q.Name = strings.Join(names, ".")

		offset++
		q.QType = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2
		q.QClass = binary.BigEndian.Uint16(data[offset : offset+2])

		questions = append(questions, &q)
		Qdcount--
	}

	return questions
}
