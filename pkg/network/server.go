package network

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
)

func StartServer(ipAndPort string) (*net.UDPConn, chan types.Probe, error) {
	fmt.Println("starting server")
	ch := make(chan types.Probe, 10)
	conn, err := StartConnection(ipAndPort)
	if err != nil {
		return nil, nil, err
	}
	go HandleProbePacket(conn, ch)
	return conn, ch, nil
}

func HandleProbePacket(conn *net.UDPConn, ch chan types.Probe) {
	fmt.Println("awaiting probes")
	for {
		probePayload := make([]byte, 12)
		n, addr, err := conn.ReadFromUDP(probePayload)
		timeNow := time.Now()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			continue
		}
		if n != 12 {
			continue
		}
		probe := types.ParseProbe(probePayload)
		probe.InitTime = timeNow
		probe.Peer = addr
		ch <- probe
	}
}
