package network

import "net"

func StartConnection(ipAndPort string) (*net.UDPConn, error) {
	addr, err := net.ResolveUDPAddr("udp", ipAndPort)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
