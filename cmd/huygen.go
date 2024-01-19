package main

import (
	"fmt"
	"os"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/controller"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
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
	c := controller.NewProbeController(conn, ch, peerIpAndPort)
	c.Run()

	time.Sleep(time.Second * 1000)
}
