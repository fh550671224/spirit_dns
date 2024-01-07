package network

import (
	"fmt"
	"net"
)

func SendUDP(data []byte, addr *net.UDPAddr) ([]byte, error) {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("net.DialUDP err: %v", err)

	}
	defer conn.Close()

	// send request
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("net.Write err: %v", err)
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("net.ReadFromUDP err: %v", err)
	}

	return buffer[:n], nil
}
