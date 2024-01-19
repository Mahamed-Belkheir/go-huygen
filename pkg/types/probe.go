package types

import (
	"encoding/binary"
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

func (pg ProbeGroup) Valid() bool {
	return time.Now().Add(time.Second).After(pg.StartTime)
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
