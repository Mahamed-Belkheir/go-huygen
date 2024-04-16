package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/lrpc"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/network"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/udpt"
	"github.com/Mahamed-Belkheir/go-huygen/pkg/variables"
)

func main() {
	addr := variables.GetListeningAddress()
	conn, err := network.StartConnection(addr)
	if err != nil {
		panic(err)
	}

	lrpcGW := lrpc.NewGateway(udpt.NewUDPTConn(conn))
	var mtx sync.Mutex
	peers := variables.GetPeers()

	file := makeFile(addr)
	offSetCh := make(chan int64, 50)
	go func(ch chan int64) {
		for {
			offSet := <-ch
			file.WriteString(fmt.Sprintf("%d,\n", offSet))
		}
	}(offSetCh)

	time.Sleep(time.Second * 2)
	for {
		for _, peer := range peers {
			synchronizeWithPeer(peer, lrpcGW, &mtx, offSetCh)
		}
		time.Sleep(time.Second * 4)
	}

	// lrpcGW.AddListener("InitSyncrhonization", func(p lrpc.Packet) *lrpc.Packet {

	// })

}

func synchronizeWithPeer(target string, lrpcGW *lrpc.Gateway, mtx *sync.Mutex, offsetCh chan int64) {
	defer func() {
		mtx.Unlock()
	}()
	fmt.Printf("starting peer synchronization with \"%s\"\n", target)
	// coded probes
	probesPacket, err := lrpcGW.RequestWithResponse(variables.GetSensorAddress(), "START_PROBE", []byte(target))
	if err != nil {
		fmt.Println("gettings probes failed", err)
		return
	}

	mtx.Lock()
	var probes map[string]uint64
	err = json.Unmarshal(probesPacket.Payload, &probes)
	if err != nil {
		fmt.Println("failed to parse probes", err)
		return
	}

	removeFirstDigit := uint64(1700000000000000000)
	t1, t2, r1, r2, pt1, pt2 := int64(probes["t1"]-removeFirstDigit), int64(probes["t2"]-removeFirstDigit), int64(probes["r1"]-removeFirstDigit), int64(probes["r2"]-removeFirstDigit), int64(probes["pt1"]-removeFirstDigit), int64(probes["pt2"]-removeFirstDigit)
	if !validCodedProbes(t1, t2, r1, r2) {
		fmt.Println("invalid probe data")
		return
	}

	newOffSet := getClockOffset(t1, t2, r1, r2, pt1, pt2)
	fmt.Printf("new offset: %d\n", newOffSet)
	offsetCh <- newOffSet
}

func validCodedProbes(t1, t2, r1, r2 int64) bool {
	return (r2-r1)-(t2-t1) < 50_000_000 && r2 > r1
}

func getClockOffset(t1, t2, r1, r2, pt1, pt2 int64) int64 {
	o1 := ((pt1 - t1) + (pt1 - r1)) / 2
	o2 := ((pt2 - t2) + (pt2 - r2)) / 2
	return (o1 + o2) / 2
}

func makeFile(ownIPAndPort string) *os.File {
	ownIPAndPort = strings.ReplaceAll(ownIPAndPort, ":", "_")
	workingDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("error creating file for collector: %s", err))
	}
	err = os.Mkdir(path.Join(workingDir, "data"), os.ModePerm)
	if err != nil && !errors.Is(err, os.ErrExist) {
		panic(fmt.Sprintf("error creating file for collector: %s", err.Error()))
	}
	f, err := os.Create(path.Join(workingDir, "data", fmt.Sprintf("%s-%d.csv", ownIPAndPort, time.Now().Unix())))
	if err != nil {
		panic(fmt.Sprintf("error creating file for collector: %s", err))
	}
	return f
}
