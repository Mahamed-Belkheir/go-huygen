package variables

import "time"

var PROBE_GROUP_DELAY = time.Second * 3
var PROBE_DELAY = time.Millisecond * 150
var MAX_PROBE_DELAY = int64(float64(PROBE_DELAY) * 1.2)
var PROBE_COUNT uint8 = 5
