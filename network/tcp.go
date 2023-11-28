package network

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
)

func sendTCPRequest(data []byte, addr *net.TCPAddr) ([]byte, error) {
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("net.DialTCP err: %s", err.Error())

	}
	defer tcpConn.Close()

	dataLength := len(data)
	tcpData := make([]byte, 2+dataLength)
	tcpData[0] = byte(dataLength >> 8)
	tcpData[1] = byte(dataLength & 0xff)
	copy(tcpData[2:], data)

	_, err = tcpConn.Write(tcpData)
	if err != nil {
		return nil, fmt.Errorf("net.Write err: %s", err.Error())

	}

	buffer := make([]byte, 4096)
	_, err = tcpConn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("net.Read err: %s", err.Error())
	}
	respLength := binary.BigEndian.Uint16(buffer[0:2])

	return buffer[2 : 2+respLength], nil
}

func TrySendTCPRequest(ipList []string, port int, data []byte) ([]byte, error) {
	for i := 0; i < len(ipList); i++ {
		ip := ipList[i]
		buffer, err := sendTCPRequest(data, &net.TCPAddr{
			IP:   net.ParseIP(ipList[i]),
			Port: port,
		})

		if err != nil {
			log.Printf("%s sendTCPRequest err: %v, trying next ip", ip, err)
		} else {
			log.Printf("%s sendTCPRequest success", ip)
			return buffer, nil
		}
	}

	return nil, fmt.Errorf("no available ip")
}
