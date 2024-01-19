package network_test

import (
	"net"
	"testing"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
)

func TestHandleConnections(t *testing.T) {

	serverConn, err := network.StartConnection("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	clientConn, err := network.StartConnection("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	ch := make(chan types.Probe)
	go network.HandleProbePacket(serverConn, ch)
	probeData := []byte{1, 0, 0, 1, 0, 96, 24, 171, 237, 239, 116, 17}
	_, err = clientConn.WriteToUDP(probeData, serverConn.LocalAddr().(*net.UDPAddr))
	if err != nil {
		panic(err)
	}

	select {
	case probe := <-ch:
		if probe.GroupId != 1 {
			t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.GroupId, 1)
		}

		if probe.Order != 1 {
			t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.Order, 1)
		}

		if probe.Type != types.SEND {
			t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.Type, types.SEND)
		}

		testTime, err := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
		if err != nil {
			panic(err)
		}
		tx := uint64(testTime.UnixNano())
		if probe.Timestamp != tx {
			t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.Timestamp, tx)
		}
		return
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for probe")
	}
}
