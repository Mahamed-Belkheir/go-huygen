package types_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/Mahamed-Belkheir/go-huygen/pkg/types"
)

func TestSerializeProbe(t *testing.T) {
	testTime, err := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	if err != nil {
		panic(err)
	}
	tx := uint64(testTime.UnixNano())
	serializedData := types.CreateSerializedProbe(1, types.SEND, 1, tx)
	expectedData := []byte{1, 0, 0, 1, 0, 96, 24, 171, 237, 239, 116, 17}

	if !bytes.Equal(serializedData, expectedData) {
		t.Fatalf("serialized data did not match expected value, got: %v, expected: %v", serializedData, expectedData)
	}
}

func TestParseProbe(t *testing.T) {
	probeData := []byte{1, 0, 0, 1, 0, 96, 24, 171, 237, 239, 116, 17}
	probe := types.ParseProbe(probeData)

	if probe.GroupId != 1 {
		t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.GroupId, 1)
	}

	if probe.Order != 1 {
		t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.Order, 1)
	}

	if probe.Type != types.SEND {
		t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.Type, types.SEND)
	}

	testTime, err := time.Parse(time.RFC3339, "2009-11-10T23:00:00Z")
	if err != nil {
		panic(err)
	}
	tx := uint64(testTime.UnixNano())
	if probe.Timestamp != tx {
		t.Fatalf("parsed data did not match expected value, got: %v, expected: %v", probe.Timestamp, tx)
	}
}
