package main

import (
	"fmt"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/lrpc"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/udpt"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/variables"
)

func main() {
	conn, err := network.StartConnection(variables.GetListeningAddress())
	if err != nil {
		panic(err)
	}

	lrpcGW := lrpc.NewGateway(udpt.NewUDPTConn(conn))

	resp, err := lrpcGW.RequestWithResponse("127.0.0.1:3001", "START_PROBE", []byte("127.0.0.1:3002"))
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(string(resp.Payload))
}
