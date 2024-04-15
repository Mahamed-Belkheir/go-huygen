package udpt

import (
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
)

type UDPTConn struct {
	latency    map[string]int
	jitter     float64
	packetLoss float64
	realConn   *net.UDPConn
}

func NewUDPTConn(conn *net.UDPConn) UDPTConn {
	return UDPTConn{
		jitter:     getFloat("HUYGENS_JITTER", 0),
		packetLoss: getFloat("HUYGENS_PACKET_LOSS", -1),
		realConn:   conn,
		latency:    getLatencyMap(),
	}
}

func getFloat(field string, defaultValue float64) float64 {
	jitterStr := os.Getenv(field)
	if jitterStr == "" {
		return defaultValue
	}
	v, err := strconv.ParseFloat(jitterStr, 64)
	if err != nil {
		panic(err)
	}
	return v
}

func getLatencyMap() map[string]int {
	latencyMap := make(map[string]int)
	latencyString := os.Getenv("HUYGENS_PEERS_LATENCY")
	if latencyString != "" {
		arr := strings.Split(latencyString, ",")
		for _, latStr := range arr {
			x := strings.Split(latStr, "=")
			latencyValue, err := strconv.ParseInt(x[1], 10, 32)
			if err != nil {
				panic(err)
			}
			latencyMap[x[0]] = int(latencyValue)
		}
	}
	return latencyMap
}

func (u *UDPTConn) ReadFromUDP(b []byte) (n int, addr *net.UDPAddr, err error) {
	return u.realConn.ReadFromUDP(b)
}

func (u *UDPTConn) WriteToUDP(b []byte, addr *net.UDPAddr) (int, error) {
	baseLatency, ok := u.latency[addr.String()]
	if ok {
		if rand.Float64() < u.packetLoss {
			return len(b), nil
		}
		latency := baseLatency + int(rand.Float64()*float64(baseLatency)*u.jitter)
		time.Sleep(time.Millisecond * time.Duration(latency))
	}
	return u.realConn.WriteToUDP(b, addr)
}
