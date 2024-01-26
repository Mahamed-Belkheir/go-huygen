package main

import (
	"fmt"
	"os"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/controller"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/variables"
)

func main() {
	ownIpAndPort := os.Args[1]
	peerIpAndPort := os.Args[2]
	fmt.Printf("ip and ports configured: own: %s, peer: %s", ownIpAndPort, peerIpAndPort)
	conn, ch, err := network.StartServer(ownIpAndPort)
	if err != nil {
		fmt.Println("failed to start server")
		panic(err)
	}
	p := controller.NewProbeTimeProcessor(variables.MAX_PROBE_DELAY)
	c := controller.NewProbeController(conn, ch, peerIpAndPort, p)
	p.Run()
	c.Run()

	select {}
}
