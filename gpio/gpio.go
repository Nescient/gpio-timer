// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/loov/hrtime"
	"github.com/warthog618/gpiod"
	"log"
	"time"
)

var gpioChip = "gpiochip2"
var startGpio = 10
var lane1Gpio = 11
var lane2Gpio = 11
var lane3Gpio = 11
var lane4Gpio = 11

var startTime time.Duration
var startCount hrtime.Count

// init will
func init() {

}

// setStartTime sets the time that the gate started
func setStartTime(evt gpiod.LineEvent) {
	gpioNum := evt.Offset // an int
	time := evt.Timestamp // time.Duration
	startTime = hrtime.Now()
	startCount = hrtime.TSC()
	log.Printf("got event %d, expecting %d\n", gpioNum, startGpio)
	log.Printf("got gate start at %v, %v, %d\n", time, startTime, startCount)
}

func Arm() (*gpiod.Line, error) {
	// gpiod.WithBothEdges and then we wont care really ?
	return gpiod.RequestLine(gpioChip, startGpio, gpiod.AsInput,
		gpiod.WithEventHandler(setStartTime), gpiod.LineEdgeRising)
}

// GetGateTime will watch the start GPIO and return a high-precision
// time value for when it starts
func GetGateTime() (time.Duration, hrtime.Count) {
	return hrtime.Now(), hrtime.TSC()
}

func handler(evt gpiod.LineEvent) {
	// handle edge event
}

func x() {
	c, err := gpiod.NewChip("gpiochip2")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(c)
	//   l, _ := c.RequestLine(rpi.J8p7, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)
	//   in, _ := gpiod.RequestLine("gpiochip0", 2, gpiod.AsInput)
	// val, _ := in.Value()
	// out, _ := gpiod.RequestLine("gpiochip0", 3, gpiod.AsOutput(val))
}
