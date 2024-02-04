package data

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
)

type ProbeCollector struct {
	ch           chan *types.ProbeGroup
	ownIpAndPort string
	f            *os.File
}

func NewProbeProbeCollector(ownIpAndPort string) *ProbeCollector {
	return &ProbeCollector{
		ch:           make(chan *types.ProbeGroup),
		ownIpAndPort: ownIpAndPort,
	}
}

func (pc *ProbeCollector) Send(pg *types.ProbeGroup) {
	pc.ch <- pg
}

func (pc *ProbeCollector) Run() {
	pc.f = makeFile(pc.ownIpAndPort)
	go func() {
		for {
			pg := <-pc.ch
			_, err := pc.f.Write(pg.ToCSVLine())
			if err != nil {
				panic("failed to write to CSV file")
			}
		}
	}()
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
