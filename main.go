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

	n, client, err := conn.ReadFromUDP(buffer)
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Printf("Received UDP request from %s", client.IP)

	// 处理数据
	dnsHeader := parseDNSHeader(buffer[:n])
	//fmt.Printf("Query: %q\n", buffer[:n])

	fmt.Printf("%x %v %x %x %x %x\n", dnsHeader.ID, dnsHeader.Flags,
		dnsHeader.Qdcount, dnsHeader.Ancount, dnsHeader.Nscount, dnsHeader.Arcount)

	questions, _ := parseDNSQuestion(buffer[:n], dnsHeader.Qdcount)

	if !dnsHeader.Flags.QR {
		for _, q := range questions {
			fmt.Println(q.Name)
			switch q.QType {
			case shared.TYPE_A:
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
	defer conn.Close()

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
	dnsHeader := parseDNSHeader(buffer[:n])
	//fmt.Printf("Response: %q\n", buffer[:n])
	if dnsHeader.Flags.TC {
		// 重新发起TCP
		for i := 0; i < len(shared.ROOT_DNS_SERVERS); i++ {
			buffer, err = tryTCP(data, shared.ROOT_DNS_SERVERS[i])
			n = len(buffer)
			if err != nil {
				if i == len(shared.ROOT_DNS_SERVERS)-1 {
					log.Fatal("no root server tcp available")
				}
				log.Printf("%s tryTCP err: %v, trying next ip", ip, err)
			} else {
				log.Printf("%s tryTCP success", ip)
				break
			}
		}
		// 重新读头
		dnsHeader = parseDNSHeader(buffer[:n])
	}

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

func tryTCP(data []byte, ip string) ([]byte, error) {
	tcpConn, err := net.DialTCP("tcp", nil, &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: 53,
	})
	if err != nil {
		return nil, err
	}
	defer tcpConn.Close()

	dataLength := len(data)
	tcpData := make([]byte, 2+dataLength)
	tcpData[0] = byte(dataLength >> 8)
	tcpData[1] = byte(dataLength & 0xff)
	copy(tcpData[2:], data)

	_, err = tcpConn.Write(tcpData)
	if err != nil {
		return nil, err
	}

	buffer := make([]byte, 4096)
	_, err = tcpConn.Read(buffer)
	if err != nil {
		return nil, err
	}
	respLength := binary.BigEndian.Uint16(buffer[0:2])

	return buffer[2 : 2+respLength], nil
}

func parseDNSHeader(data []byte) *DNSHeader {
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

func checkIsPointer(data []byte, offset int) (bool, uint16) {
	if data[offset]&0xC0 == 0xC0 {
		ptr := binary.BigEndian.Uint16(data[offset:offset+2]) - 0xC000
		return true, ptr
	} else {
		return false, 0
	}
}

func parseNameByOffset(data []byte, offset int) (string, int) {
	var names []string
	for {
		if data[offset] == 0 {
			offset++
			break
		}

		if isPointer, ptr := checkIsPointer(data, offset); isPointer {
			name, _ := parseNameByOffset(data, int(ptr))
			names = append(names, name)
			offset += 2
			if data[offset] == 0 {
				break
			}
		} else {
			length := int(data[offset])
			offset++
			name := data[offset : offset+length]
			names = append(names, string(name))
			offset += length
		}
	}

	return strings.Join(names, "."), offset
}

func parseDNSResourceRecord(data []byte, offset int, count uint16) ([]*DNSResourceRecord, int) {
	var records []*DNSResourceRecord
	for count > 0 {
		r := DNSResourceRecord{}

		r.Name, offset = parseNameByOffset(data, offset)

		r.RType = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.RClass = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		r.TTL = binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4

		r.ResourceDataLength = binary.BigEndian.Uint16(data[offset : offset+2])
		offset += 2

		switch r.RType {
		case shared.TYPE_A:
			r.ResourceData = parseARecordData(data, offset, r.ResourceDataLength)
		case shared.TYPE_NS, shared.TYPE_CNAME:
			r.ResourceData, _ = parseNameByOffset(data, offset)
		}
		offset += int(r.ResourceDataLength)

		records = append(records, &r)
		count--
	}

	return records, offset
}

func parseARecordData(data []byte, offset int, length uint16) string {
	var parts []string
	for length > 0 {
		parts = append(parts, fmt.Sprintf("%d", data[offset]))
		offset++
		length--
	}
	return strings.Join(parts, ".")
}
