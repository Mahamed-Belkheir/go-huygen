package lrpc_test

import (
	"fmt"
	"testing"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/lrpc"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
)

// func TestSerialization(t *testing.T) {
// 	serverConn, err := network.StartConnection("127.0.0.1:8083")
// 	if err != nil {
// 		panic(err)
// 	}
// 	serverGateway := lrpc.NewGateway(*serverConn)

// 	msg, err := serverGateway.ConstructMessage("function", []byte("test data"), nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(msg)

// 	p := lrpc.ParsePayload(msg, &net.UDPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 3033})
// 	fmt.Println(string(p.Payload))
// }

func TestLrpc(t *testing.T) {
	serverConn, err := network.StartConnection("127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	clientConn, err := network.StartConnection("127.0.0.1:8081")
	if err != nil {
		panic(err)
	}

	serverGateway := lrpc.NewGateway(*serverConn)
	clientGateway := lrpc.NewGateway(*clientConn)

	serverGateway.AddListener("TestFunc", func(p lrpc.Packet) *lrpc.Packet {
		fmt.Println(p.Source)
		p.Payload = []byte("return hello")
		return &p
	})

	response, err := clientGateway.RequestWithResponse("127.0.0.1:8080", "TestFunc", []byte("Hello world"))

	fmt.Println(response, err)
}
