package controller

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
)

type ProbeController struct {
	trackedProbes map[uint16]*types.ProbeGroup
	c             network.Client
	ch            chan types.Probe
	peer          string
	mtx           *sync.Mutex
}

func NewProbeController(conn *net.UDPConn, ch chan types.Probe, peer string) ProbeController {
	return ProbeController{
		trackedProbes: map[uint16]*types.ProbeGroup{},
		c:             network.NewClient(conn),
		ch:            ch,
		peer:          peer,
		mtx:           &sync.Mutex{},
	}
}

func (p ProbeController) Run() {
	go p.clientProbing()
	go p.handleServerMessages()
}

func (p ProbeController) clientProbing() {
	var groupId uint16 = 0
probeGroup:
	for {
		time.Sleep(time.Second * 5)
		groupId += 1
		p.mtx.Lock()
		fmt.Println("starting new group", groupId)
		p.trackedProbes[groupId] = &types.ProbeGroup{
			Id: groupId,
		}
		var i uint8

		for i = 0; i < 5; i++ {
			timeNow := time.Now()
			ts := uint64(timeNow.UnixNano())
			payload := types.CreateSerializedProbe(groupId, types.SEND, i, ts)
			err := p.c.Send(p.peer, payload)
			if err != nil {
				fmt.Println("error sending probe in group: ", groupId, err)
				continue probeGroup
			}
			probe := types.ParseProbe(payload)
			probe.InitTime = timeNow
			pg := p.trackedProbes[groupId]
			pg.Sent = append(pg.Sent, probe)
			p.trackedProbes[groupId] = pg
			time.Sleep(time.Millisecond * 150)
		}
		p.mtx.Unlock()
	}
}

func (p ProbeController) handleServerMessages() {
	for {
		probe := <-p.ch
		if probe.Type == types.RECV {
			p.mtx.Lock()
			pg, ok := p.trackedProbes[probe.GroupId]
			if !ok {
				fmt.Println("probe group not found:", probe.GroupId)
				continue
			}
			pg.Received = append(pg.Received, probe)
			if probe.Order == 4 {
				if pg.AllProbesInOrder() {
					pg.PrintDeltas()
				} else {
					fmt.Println("got invalid probe group: ", probe.GroupId)
				}
			}
			p.mtx.Unlock()
		} else {
			ts := uint64(time.Now().UnixNano())
			payload := types.CreateSerializedProbe(probe.GroupId, types.RECV, probe.Order, ts)
			p.c.Send(p.peer, payload)
		}
	}
}
