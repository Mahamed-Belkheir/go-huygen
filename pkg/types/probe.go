package types

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

type probeType = uint8

var (
	SEND probeType = 0
	RECV probeType = 1
)

/*
Probe type, can be serialized and parsed, the serialized form is 12 bytes long

2 bytes group id

1 byte type

1 byte order

8 bytes timestamp
*/
type Probe struct {
	GroupId   uint16
	Type      probeType
	Order     uint8
	Timestamp uint64
	InitTime  time.Time
	Peer      *net.UDPAddr
}

type ProbeGroup struct {
	Id        uint16
	Received  []Probe
	Sent      []Probe
	StartTime time.Time
}

func (pg ProbeGroup) AllProbesInOrder() bool {
	if len(pg.Received) != len(pg.Sent) {
		fmt.Println("not equal length:", len(pg.Received), len(pg.Sent))
		return false
	}
	for i, probe := range pg.Sent {
		expectedOrder := uint8(i)
		if probe.Order != expectedOrder || pg.Received[i].Order != expectedOrder {
			fmt.Println("bad order:", probe.Order, pg.Received[i].Order)
			return false
		}
	}
	return true
}

func (pg ProbeGroup) PrintDeltas() {
	for i := range pg.Sent[:len(pg.Sent)-1] {
		fmt.Printf("probe : %d \n", i)
		fmt.Printf("round trip delay: %dus\n", (int64(pg.Received[i].InitTime.UnixNano())-int64(pg.Sent[i].InitTime.UnixNano()))/1000)
		fmt.Printf("%d to %d delay: %dms\n", i, i+1, (int64(pg.Received[i+1].InitTime.UnixNano())-int64(pg.Received[i].InitTime.UnixNano()))/1000_000)
	}
	fmt.Print("\n\n")
}

func (pg ProbeGroup) DeltaBetweenProbes() []int64 {
	deltaBetweenProbes := []int64{}
	for i := range pg.Sent[:len(pg.Sent)-1] {
		deltaBetweenProbes = append(deltaBetweenProbes, (int64(pg.Received[i+1].InitTime.UnixNano()) - int64(pg.Received[i].InitTime.UnixNano())))
	}
	return deltaBetweenProbes
}

func (pg ProbeGroup) AverageOWDEstimateOfGroup() int64 {
	var total int64 = 0
	for i := range pg.Sent {
		total += (int64(pg.Received[i].InitTime.UnixNano()) - int64(pg.Sent[i].InitTime.UnixNano())) / 2
	}
	return total / int64(len(pg.Received))
}

func ParseProbe(rawData []byte) Probe {
	return Probe{
		GroupId:   binary.LittleEndian.Uint16(rawData[0:2]),
		Type:      probeType(rawData[2]),
		Order:     rawData[3],
		Timestamp: binary.LittleEndian.Uint64(rawData[4:]),
	}
}

func CreateSerializedProbe(groupId uint16, pType, order uint8, timestamp uint64) []byte {
	data := []byte{}
	data = binary.LittleEndian.AppendUint16(data, groupId)
	data = append(data, pType, order)
	data = binary.LittleEndian.AppendUint64(data, timestamp)
	return data
}
