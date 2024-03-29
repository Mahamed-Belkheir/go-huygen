package controller

import (
	"fmt"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/data"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
)

type ProbeTimeStats struct {
	validGroups     int
	invalidGroups   int
	lastOWDEstimate int64
}

type ProbeTimeProcessor struct {
	parsedGroups          chan *types.ProbeGroup
	started               bool
	maxDelayBetweenProbes int64
	stats                 ProbeTimeStats
	collector             *data.ProbeCollector
}

func NewProbeTimeProcessor(maxDelayBetweenProbes int64, collector *data.ProbeCollector) *ProbeTimeProcessor {
	return &ProbeTimeProcessor{
		parsedGroups:          make(chan *types.ProbeGroup, 1024),
		started:               false,
		maxDelayBetweenProbes: maxDelayBetweenProbes,
		stats: ProbeTimeStats{
			validGroups:     0,
			invalidGroups:   0,
			lastOWDEstimate: 0,
		},
		collector: collector,
	}
}

func (p ProbeTimeProcessor) AddGroup(pg *types.ProbeGroup) {
	p.parsedGroups <- pg
}

func (p ProbeTimeProcessor) Run() {
	if p.started {
		return
	}
	p.started = true
	go p.processProbeGroups()
}

func (p ProbeTimeProcessor) processProbeGroups() {
processLoop:
	for {
		pg := <-p.parsedGroups
		if !pg.AllProbesInOrder() {
			p.stats.invalidGroups += 1
			fmt.Println("got invalid group order")
			continue processLoop
		}
		deltas := pg.DeltaBetweenProbes()
		for _, d := range deltas {
			if d > p.maxDelayBetweenProbes {
				p.stats.invalidGroups += 1
				fmt.Println("got invalid group delay", d, p.maxDelayBetweenProbes)
				continue processLoop
			}
		}
		p.stats.validGroups += 1
		p.stats.lastOWDEstimate = pg.AverageOWDEstimateOfGroup()
		fmt.Printf("new OWD estimate %dus \n", p.stats.lastOWDEstimate/1000)
		p.collector.Send(pg)
	}
}
