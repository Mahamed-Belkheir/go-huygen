package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/lrpc"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/timet"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/udpt"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/variables"
)

func main() {
	conn, err := network.StartConnection(variables.GetListeningAddress())
	if err != nil {
		panic(err)
	}

	lrpcGW := lrpc.NewGateway(udpt.NewUDPTConn(conn))

	lrpcGW.AddListener("PROBE", func(p lrpc.Packet) *lrpc.Packet {
		t := timet.GetTime()
		data := []byte{}
		data = binary.LittleEndian.AppendUint64(data, t)
		p.Payload = data
		copy := p
		return &copy
	})

	lrpcGW.AddListener("START_PROBE", func(p lrpc.Packet) *lrpc.Packet {
		target := string(p.Payload)
		fmt.Println("starting probe", target)
		var wg sync.WaitGroup
		wg.Add(2)

		var t1, t2, r1, r2, pt1, pt2 uint64
		var err1, err2 error

		go func() {
			t1, r1, pt1, err1 = callProbe(target, lrpcGW)
			wg.Done()
		}()

		time.Sleep(time.Millisecond * 50)

		go func() {
			t2, r2, pt2, err2 = callProbe(target, lrpcGW)
			wg.Done()
		}()

		wg.Wait()

		if err1 != nil || err2 != nil {
			fmt.Println(err)
			return nil
		}

		data, err := json.Marshal(map[string]uint64{
			"t1":  t1,
			"t2":  t2,
			"r1":  r1,
			"r2":  r2,
			"pt1": pt1,
			"pt2": pt2,
		})

		if err != nil {
			fmt.Println(err)
			return nil
		}
		p.Payload = data
		fmt.Println("done probing")
		return &p
	})

	select {}
}

func callProbe(target string, lrpcGW lrpc.Gateway) (uint64, uint64, uint64, error) {
	t := timet.GetTime()
	response, err := lrpcGW.RequestWithResponse(target, "PROBE", []byte{})

	r := timet.GetTime()
	if err != nil {
		return 0, 0, 0, err
	}
	pt := binary.LittleEndian.Uint64(response.Payload)
	return t, r, pt, nil
}
