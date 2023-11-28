package network

import (
	"fmt"
	"log"
	"net"
)

func sendUDPRequest(data []byte, addr *net.UDPAddr) ([]byte, error) {
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("net.DialUDP err: %s", err.Error())

	}
	defer conn.Close()

	// send request
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("net.Write err: %s", err.Error())
	}

	buffer := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, fmt.Errorf("net.ReadFromUDP err: %s", err.Error())
	}

	return buffer[:n], nil
}

func TrySendUDPRequest(ipList []string, port int, data []byte) ([]byte, error) {
	for i := 0; i < len(ipList); i++ {
		ip := ipList[i]
		buffer, err := sendUDPRequest(data, &net.UDPAddr{
			IP:   net.ParseIP(ipList[i]),
			Port: port,
		})

		if err != nil {
			log.Printf("%s sendUDPRequest err: %v, trying next ip", ip, err)
		} else {
			log.Printf("%s sendUDPRequest success", ip)
			return buffer, nil
		}
	}
	return nil, fmt.Errorf("no available ip")
}
