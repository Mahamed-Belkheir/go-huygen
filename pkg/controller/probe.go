package controller

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/variables"
)

type ProbeController struct {
	trackedProbes    map[uint16]*types.ProbeGroup
	c                network.Client
	probesFromServer chan types.Probe
	validateProbes   chan types.Probe
	timePorcessor    *ProbeTimeProcessor
	peer             string
	mtx              *sync.Mutex
}

func NewProbeController(conn *net.UDPConn, ch chan types.Probe, peer string, timeProcessor *ProbeTimeProcessor) ProbeController {
	return ProbeController{
		trackedProbes:    map[uint16]*types.ProbeGroup{},
		c:                network.NewClient(conn),
		probesFromServer: ch,
		validateProbes:   make(chan types.Probe, 100),
		peer:             peer,
		mtx:              &sync.Mutex{},
		timePorcessor:    timeProcessor,
	}
}

func (p ProbeController) Run() {
	go p.handleProbes()
	time.Sleep(time.Millisecond)
	go p.clientProbing()
	time.Sleep(time.Millisecond)
	go p.handleServerMessages()
}

func (p ProbeController) handleProbes() {
	for {
		probe := <-p.validateProbes
		pg, ok := p.trackedProbes[probe.GroupId]
		if probe.Type == types.SEND {
			if !ok {
				fmt.Println("setup group")
				pg = &types.ProbeGroup{
					Id:        probe.GroupId,
					StartTime: probe.InitTime,
				}
				p.trackedProbes[probe.GroupId] = pg
			}
			pg.Sent = append(pg.Sent, probe)
		} else {
			if !ok || pg == nil {
				fmt.Println("group not found")
				continue
			}
			pg.Received = append(pg.Received, probe)

			if probe.Order == uint8(variables.PROBE_COUNT) {
				p.timePorcessor.AddGroup(pg)
			}
		}
	}
}

func (p ProbeController) clientProbing() {
	var groupId uint16 = 0
probeGroup:
	for {
		time.Sleep(variables.PROBE_GROUP_DELAY)
		groupId += 1
		fmt.Println("starting new group", groupId)
		var i uint8
		for i = 0; i <= variables.PROBE_COUNT; i++ {
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
			p.validateProbes <- probe
			time.Sleep(variables.PROBE_DELAY)
		}
	}
}

func (p ProbeController) handleServerMessages() {
	for {
		probe := <-p.probesFromServer
		if probe.Type == types.RECV {
			p.validateProbes <- probe
		} else {
			ts := uint64(time.Now().UnixNano())
			payload := types.CreateSerializedProbe(probe.GroupId, types.RECV, probe.Order, ts)
			p.c.Send(p.peer, payload)
		}
	}
}
