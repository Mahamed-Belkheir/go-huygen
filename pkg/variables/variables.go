package variables

import "time"

var PROBE_GROUP_DELAY = time.Second * 5
var PROBE_DELAY = time.Millisecond * 150
var MAX_PROBE_DELAY = int64(float64(PROBE_DELAY) * 1.1)
var PROBE_COUNT uint8 = 10
