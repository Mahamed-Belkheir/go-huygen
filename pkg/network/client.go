package network

import (
	"net"
)

type Client struct {
	conn *net.UDPConn
}

func NewClient(conn *net.UDPConn) Client {
	return Client{conn}
}

func (c Client) Send(ipAndPort string, data []byte) error {
	addr, err := net.ResolveUDPAddr("udp", ipAndPort)
	if err != nil {
		return err
	}
	_, err = c.conn.WriteToUDP(data, addr)
	if err != nil {
		return err
	}
	return nil
}
