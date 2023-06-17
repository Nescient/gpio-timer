// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/warthog618/gpiod"
	"github.com/loov/hrtime"
	"time"
	"log"
)

var startGpio = 10
var lane1Gpio = 11
var lane2Gpio = 11
var lane3Gpio = 11
var lane4Gpio = 11


// init will 
func init() {
	
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
	if err != nil{
		log.Fatal(err)
	}
	log.Println(c);
//   l, _ := c.RequestLine(rpi.J8p7, gpiod.WithEventHandler(handler), gpiod.WithBothEdges)
//   in, _ := gpiod.RequestLine("gpiochip0", 2, gpiod.AsInput)
// val, _ := in.Value()
// out, _ := gpiod.RequestLine("gpiochip0", 3, gpiod.AsOutput(val))
  }