package lrpc

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"time"
)

const messageSize = 64000

type Packet struct {
	ID      uint16
	Source  string
	Command string
	Payload []byte
}

type Gateway struct {
	id                uint16
	udpConn           net.UDPConn
	responseListeners map[uint16]chan Packet
	commandListeners  map[string]chan Packet
}

func NewGateway(udpConn net.UDPConn) Gateway {
	g := Gateway{
		0,
		udpConn,
		make(map[uint16]chan Packet),
		make(map[string]chan Packet),
	}
	go g.listenForMessages()
	return g
}

func (g *Gateway) Request(target string, name string, payload []byte) (uint16, error) {
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return 0, err
	}
	data, err := g.ConstructMessage(name, payload, nil)
	if err != nil {
		return 0, err
	}
	g.udpConn.WriteToUDP(data, addr)
	return g.id, nil
}

func (g *Gateway) RequestWithResponse(target string, name string, payload []byte) (Packet, error) {
	id, err := g.Request(target, name, payload)
	if err != nil {
		return Packet{}, err
	}
	packetChan := make(chan Packet, 10)
	g.responseListeners[id] = packetChan
	select {
	case response := <-packetChan:
		return response, nil
	case <-time.After(time.Second * 5):
		return Packet{}, errors.New("timed out")
	}
}

func (g *Gateway) AddListener(name string, cb func(Packet) *Packet) {
	ch := make(chan Packet, 10)
	go func() {
		for {
			packet := <-ch
			response := cb(packet)
			if response != nil {
				g.SendResponse(packet)
			}
		}
	}()
	g.commandListeners[name] = ch
}

func (g *Gateway) SendResponse(packet Packet) error {
	fmt.Println("sending response")
	addr, err := net.ResolveUDPAddr("udp", packet.Source)
	if err != nil {
		return err
	}
	data, err := g.ConstructMessage("Response", packet.Payload, &packet.ID)
	if err != nil {
		return err
	}
	g.udpConn.WriteToUDP(data, addr)
	return nil
}

func (g *Gateway) ConstructMessage(name string, payload []byte, idInput *uint16) ([]byte, error) {
	data := []byte{}
	var id uint16
	if idInput == nil {
		g.id = g.id + 1
		id = g.id
	} else {
		id = *idInput
	}
	data = binary.LittleEndian.AppendUint16(data, id)
	data = binary.LittleEndian.AppendUint16(data, uint16(len(name)))
	data = append(data, []byte(name)...)
	data = binary.LittleEndian.AppendUint16(data, uint16(len(payload)))
	data = append(data, payload...)
	return data, nil
}

func (g *Gateway) listenForMessages() {
	messageBuffer := make([]byte, messageSize)
	for {
		_, addr, err := g.udpConn.ReadFromUDP(messageBuffer)
		if err != nil {
			continue
		}
		p := ParsePayload(messageBuffer, addr)
		if p.Command == "Response" {
			fmt.Println("got response", p.ID, g.responseListeners)
			g.responseListeners[p.ID] <- p
		} else {
			ch, ok := g.commandListeners[p.Command]
			if !ok {
				continue
			}
			fmt.Println("got ch")
			ch <- p
		}
	}
}

func ParsePayload(rawData []byte, addr *net.UDPAddr) Packet {
	commandSize := binary.LittleEndian.Uint16(rawData[2:4])
	payloadSize := binary.LittleEndian.Uint16(rawData[commandSize+4 : commandSize+7])
	return Packet{
		ID:      binary.LittleEndian.Uint16(rawData[0:2]),
		Command: string(rawData[4 : 4+commandSize]),
		Source:  addr.String(),
		Payload: rawData[commandSize+6 : commandSize+6+payloadSize],
	}
}
