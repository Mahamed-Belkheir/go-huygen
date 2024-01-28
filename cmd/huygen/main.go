package main

import (
	"fmt"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/controller"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/variables"
)

func main() {
	ownIpAndPort := variables.GetListeningAddress()
	peers := variables.GetPeers()

	fmt.Println("starting server at:", ownIpAndPort)
	conn, mainInbound, err := network.StartServer(ownIpAndPort)
	if err != nil {
		fmt.Println("failed to start server")
		panic(err)
	}

	peerChannels := make(map[string]chan types.Probe, 1024)

	for _, peerIpAndPort := range peers {
		ch := make(chan types.Probe)
		peerChannels[peerIpAndPort] = ch
		go func(peerIpAndPort string) {
			fmt.Println("setting up peer at:", peerIpAndPort)
			p := controller.NewProbeTimeProcessor(variables.MAX_PROBE_DELAY)
			c := controller.NewProbeController(conn, ch, peerIpAndPort, p)
			p.Run()
			c.Run()
		}(peerIpAndPort)
	}

	go func() {
		for {
			probe := <-mainInbound
			peerChannels[probe.Peer.String()] <- probe
		}
	}()
	select {}
}
