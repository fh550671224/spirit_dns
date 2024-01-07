package network

import (
	"encoding/binary"
	"fmt"
	"net"
)

func SendTCPRequest(data []byte, addr *net.TCPAddr) ([]byte, error) {
	tcpConn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		return nil, fmt.Errorf("net.DialTCP err: %v", err)

	}
	defer tcpConn.Close()

	dataLength := len(data)
	tcpData := make([]byte, 2+dataLength)
	tcpData[0] = byte(dataLength >> 8)
	tcpData[1] = byte(dataLength & 0xff)
	copy(tcpData[2:], data)

	_, err = tcpConn.Write(tcpData)
	if err != nil {
		return nil, fmt.Errorf("net.Write err: %v", err)

	}

	buffer := make([]byte, 4096)
	_, err = tcpConn.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("net.Read err: %v", err)
	}
	respLength := binary.BigEndian.Uint16(buffer[0:2])

	return buffer[2 : 2+respLength], nil
}
