// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/warthog618/gpiod"
	"log"
	"sync"
	"time"
)

// pin mapping, see https://hub.libre.computer/t/libre-computer-wiring-tool/40
// Pin     Chip    Line    sysfs   Name    Pad     Ref     Desc
// 1       3.3V    3.3V    3.3V    3.3V    3.3V    VCC_IO  3.3V
// 3       2       25      89      GPIO2_D1        R17     
// 5       2       24      88      GPIO2_D0        P17     
// 6       GND     GND     GND     GND     GND     GND     GND
// 7       1       28      60      GPIO1_D4        Y17     PULL DOWN, NOT UP!
// 8       3       4       100     GPIO3_A4        E2      
// 9       GND     GND     GND     GND     GND     GND     GND
// 11      2       20      84      GPIO2_C4        V18     
// 13      2       21      85      GPIO2_C5        V17     
// 15      2       22      86      GPIO2_C6        V16     
// 17      GND     GND     GND     GND     GND     VCC_IO  GND

var startChip = "gpiochip3"
var laneChip = "gpiochip2"
var startGpio = 4
var lane1Gpio = 24
var lane2Gpio = 25
var lane3Gpio = 20
var lane4Gpio = 21

var startTime time.Duration
var laneTimes = [4]time.Duration{}

var waitStart sync.WaitGroup
var waitLanes sync.WaitGroup

// clearLanes resets all the lane times to a default value
func clearLanes() {
	laneTimes = [4]time.Duration{0, 0, 0, 0}
}

// setStartTime sets the time that the gate started
func setStartTime(evt gpiod.LineEvent) {
	startTime = evt.Timestamp
	waitStart.Done()
}

// setLaneTime sets the time that a given lane completes
func setLaneTime(evt gpiod.LineEvent) {
	log.Printf("got lane event %d\n", evt.Offset)
	switch gpioNum := evt.Offset; gpioNum {
	case lane1Gpio:
		if laneTimes[0] == 0 {
			laneTimes[0] = evt.Timestamp
			waitLanes.Done()
		}
	case lane2Gpio:
		if laneTimes[1] == 0 {
			laneTimes[1] = evt.Timestamp
			waitLanes.Done()
		}
	case lane3Gpio:
		if laneTimes[2] == 0 {
			laneTimes[2] = evt.Timestamp
			waitLanes.Done()
		}
	case lane4Gpio:
		if laneTimes[3] == 0 {
			laneTimes[3] = evt.Timestamp
			waitLanes.Done()
		}
	default:
		log.Printf("unknown lane event %d\n", gpioNum)
	}
}

// ArmStart sets up the interrupt handler for the start GPIO line
func ArmStart() (*gpiod.Line, error) {
	clearLanes()
	waitStart.Add(1)
	return gpiod.RequestLine(startChip, startGpio, gpiod.AsInput,
		gpiod.WithEventHandler(setStartTime), gpiod.LineEdgeFalling)
}

// ArmLanes sets up the interrupt handler for the all the lane GPIO lines
func ArmLanes() (*gpiod.Lines, error) {
	clearLanes()
	waitLanes.Add(4)
	return gpiod.RequestLines(laneChip, []int{lane1Gpio, lane2Gpio, lane3Gpio, lane4Gpio},
		gpiod.AsInput, gpiod.WithEventHandler(setLaneTime), gpiod.LineEdgeFalling)
}

// WaitForStart waits until the start GPIO triggers
func WaitForStart() {
	waitStart.Wait()
}

// deltaTimes calculates the difference betwene two timestamps
func deltaTimes(start time.Duration, end time.Duration) float64 {
	if end < start {
		return 0.0
	}
	return end.Seconds() - start.Seconds()
}

// WaitForLanes waits until all 4 lanes have triggered and returns
// the time difference for each lane
func WaitForLanes() (times [4]float64) {
	waitLanes.Wait()
	for i, _ := range laneTimes {
		times[i] = deltaTimes(startTime, laneTimes[i])
	}
	return
}
