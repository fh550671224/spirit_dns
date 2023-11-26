package spiritDNS

import (
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

	n, addr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return
	}

	// 处理数据
	dnsRequest, parseErr := parseDNSRequest(buffer[:n])
	if parseErr != nil {
		return
	}

	fmt.Println(dnsRequest, addr)

}

func parseDNSRequest(data []byte) (DNSRequest, error) {
	return DNSRequest{}, nil
}
