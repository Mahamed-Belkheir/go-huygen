package timet

import (
	"os"
	"strconv"
	"time"
)

func init() {
	offsetStr := os.Getenv("HUYGENS_OFFSET")
	if offsetStr != "" {
		offset, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			panic(err)
		}
		currentOffset = int(offset)
	}
}

var currentOffset = 0

func GetTime() uint64 {
	return uint64(time.Now().UnixNano() + int64(currentOffset))
}

func EditOffset(newOffset int) {
	currentOffset += newOffset
}
