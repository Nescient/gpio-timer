// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/warthog618/gpiod"
	"github.com/loov/hrtime"
)

var startGpio = 10
var lane1Gpio = 11
var lane2Gpio = 11
var lane3Gpio = 11
var lane4Gpio = 11

// GetGateTime will watch the start GPIO and return a high-precision
// time value for when it starts
func GetGateTime() time.Duration, int64 {
	return hrtime.Now(), hrtime.TSC()
}