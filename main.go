package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"spiritDNS/shared"
	"strings"
)

type DNSQuestion struct {
	Name   string
	QType  uint16
	QClass uint16
}

type DNSHeader struct {
	ID      uint16
	Flags   DNSFlags
	Qdcount uint16
	Ancount uint16
	Nscount uint16
	Arcount uint16
}

type DNSFlags struct {
	QR       bool  // 是否是response
	OpCode   uint8 // 操作码
	AA       bool  // 权威答案
	TC       bool  // 是否截断
	RD       bool  // 是否期望递归
	RA       bool  // 是否支持递归
	Reserved uint8 // 保留，通常为0
	RespCode uint8 // 响应码
}

type DNSResourceRecord struct {
	Name               string
	RType              uint16
	RClass             uint16
	TTL                uint32
	ResourceDataLength uint16
	ResourceData       string
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
		log.Fatal(err)
		return
	}

	// 处理数据
	dnsHeader := parseDNSRequest(buffer[:n])
	//fmt.Printf("Query: %q\n", buffer[:n])

	fmt.Printf("%x %v %x %x %x %x\n", dnsHeader.ID, dnsHeader.Flags,
		dnsHeader.Qdcount, dnsHeader.Ancount, dnsHeader.Nscount, dnsHeader.Arcount)

	questions, _ := parseDNSQuestion(buffer[:n], dnsHeader.Qdcount)

	if !dnsHeader.Flags.QR {
		for _, q := range questions {
			fmt.Println(q.Name)
			switch q.QType {
			case shared.QTYPE_A:
				fmt.Println("A record")
				// establish conn to root server
				sendDNSRequest(buffer[:n], shared.ROOT_DNS_SERVERS[0])
			}
		}
	}

	//fmt.Printf("length of last name:%x and is: %c %c %c\n", buffer[12], buffer[13], buffer[14], buffer[15])
}

func sendDNSRequest(data []byte, ip string) {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(ip),
		Port: 53,
	})
	if err != nil {
		log.Fatal(err)
	}

	// send request
	_, err = conn.Write(data)
	if err != nil {
		log.Fatal(err)
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
	}

	// 处理数据
	dnsHeader := parseDNSRequest(buffer[:n])
	fmt.Printf("Response: %q\n", buffer[:n])

	fmt.Printf("%x %v %x %x %x %x\n", dnsHeader.ID, dnsHeader.Flags,
		dnsHeader.Qdcount, dnsHeader.Ancount, dnsHeader.Nscount, dnsHeader.Arcount)

	_, offset := parseDNSQuestion(buffer[:n], dnsHeader.Qdcount)

	answerRRs, offset := parseDNSResourceRecord(buffer[:n], offset, dnsHeader.Ancount)
	fmt.Println(answerRRs)

	NSRRs, offset := parseDNSResourceRecord(buffer[:n], offset, dnsHeader.Nscount)
	fmt.Println(NSRRs)

	additionalRRs, offset := parseDNSResourceRecord(buffer[:n], offset, dnsHeader.Arcount)
	fmt.Println(additionalRRs)
}

func parseNameByOffset(data []byte, offset int) (string, int) {
	var names []string
	for data[offset] != 0 {
		length := int(data[offset])
		offset++
		name := data[offset : offset+length]
		names = append(names, string(name))
		offset += length
	}
	res := strings.Join(names, ".")
	offset++
	return res, offset
}

func parseDNSRequest(data []byte) *DNSHeader {
	return &DNSHeader{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		Flags:   *parseDNSFlags(binary.BigEndian.Uint16(data[2:4])),
		Qdcount: binary.BigEndian.Uint16(data[4:6]),
		Ancount: binary.BigEndian.Uint16(data[6:8]),
		Nscount: binary.BigEndian.Uint16(data[8:10]),
		Arcount: binary.BigEndian.Uint16(data[10:12]),
	}
}

func parseDNSFlags(flags uint16) *DNSFlags {
	return &DNSFlags{
		QR:       (flags & 0x8000) != 0,
		OpCode:   uint8((flags & 0x7800) >> 11),
		AA:       (flags & 0x0400) != 0,
		TC:       (flags & 0x0200) != 0,
		RD:       (flags & 0x0100) != 0,
		RA:       (flags & 0x0080) != 0,
		Reserved: uint8((flags & 0x0070) >> 4),
		RespCode: uint8(flags & 0x000f),
	}
}

func parseDNSQuestion(data []byte, Qdcount uint16) ([]*DNSQuestion, int) {
	offset := 12
	var questions []*DNSQuestion
	for Qdcount > 0 {
		q := DNSQuestion{}
		q.Name, offset = parseNameByOffset(data, offset)

		q.QType = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		q.QClass = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		questions = append(questions, &q)
		Qdcount--
	}

	return questions, offset
}

func parseDNSResourceRecord(data []byte, offset int, count uint16) ([]*DNSResourceRecord, int) {
	var records []*DNSResourceRecord
	for count > 0 {
		r := DNSResourceRecord{}

		if data[offset]&0xC0 != 0xC0 {
			r.Name, offset = parseNameByOffset(data, offset)

		} else {
			ptr := binary.BigEndian.Uint16(data[offset:offset+2]) - 0xC000
			r.Name, _ = parseNameByOffset(data, int(ptr))
			offset += 2
		}

		r.RType = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.RClass = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.TTL = binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4

		r.ResourceDataLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		switch r.RType {
		case shared.QTYPE_A:
			r.ResourceData = string(data[offset : offset+int(r.ResourceDataLength)])
		case shared.QTYPE_NS, shared.QTYPE_CNAME:
			r.ResourceData, _ = parseNameByOffset(data, offset)
		}
		offset += int(r.ResourceDataLength)

		records = append(records, &r)
		count--
	}

	return records, offset
}
