package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
)

type DNSRequest struct {
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

	fmt.Printf("%x %x %x %x %x %x\n", dnsHeader.ID, dnsHeader.Flags, dnsHeader.Qdcount, dnsHeader.Ancount, dnsHeader.Nscount,
		dnsHeader.Arcount)
}

func parseDNSRequest(data []byte) DNSHeader {
	return DNSHeader{
		ID:      binary.BigEndian.Uint16(data[0:2]),
		Flags:   binary.BigEndian.Uint16(data[2:4]),
		Qdcount: binary.BigEndian.Uint16(data[4:6]),
		Ancount: binary.BigEndian.Uint16(data[6:8]),
		Nscount: binary.BigEndian.Uint16(data[8:10]),
		Arcount: binary.BigEndian.Uint16(data[10:12]),
	}
}

type DNSHeader struct {
	ID      uint16
	Flags   uint16
	Qdcount uint16
	Ancount uint16
	Nscount uint16
	Arcount uint16
}
