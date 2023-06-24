// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/warthog618/gpiod"
	"log"
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

// Constants to customize this timer.
var startChip = "gpiochip3"
var laneChip = "gpiochip2"
var startGpio = 4
var lane1Gpio = 24
var lane2Gpio = 25
var lane3Gpio = 20
var lane4Gpio = 21

// GpioTime is a structure that represents a single GPIO event
type GpioTime struct {
	Chip    string
	Lane    int
	Time    time.Duration
	Pending bool
	Line    *gpiod.Line
	Channel chan int
}

// New initializes the structure to default values
func (this GpioTime) New(chip string, offset int) {
	this.Close()
	this = GpioTime{chip, offset, 0, true, nil, make(chan int)}
}

// Arm will register the GPIO for a falling edge event
func (this GpioTime) Arm() (err error) {
	this.Pending = true
	this.Line, err = gpiod.RequestLine(this.Chip, this.Lane, gpiod.AsInput,
		gpiod.WithEventHandler(this.gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	return
}

// gpioHandler handles a GPIO event for a given GpioTime struct
func (this GpioTime) gpioHandler(evt gpiod.LineEvent) {
	if evt.Offset == this.Lane {
		this.Pending = false
		this.Time = evt.Timestamp
		this.Channel <- 1
	} else {
		log.Printf("Received unknown GPIO event %d\n", evt.Offset)
	}
}

// Close will close any open GPIO lanes for the GpioTime struct
// as well as the channel
func (this GpioTime) Close() {
	if this.Line != nil {
		this.Line.Close()
	}
	close(this.Channel)
}

// createLanes initializes an array of GpioTime structures
// to represent the set of Gpio lanes
func createLanes() (lanes [4]GpioTime) {
	lanes[0] = GpioTime{laneChip, lane1Gpio, 0, true, nil, make(chan int)}
	lanes[1] = GpioTime{laneChip, lane2Gpio, 0, true, nil, make(chan int)}
	lanes[2] = GpioTime{laneChip, lane3Gpio, 0, true, nil, make(chan int)}
	lanes[3] = GpioTime{laneChip, lane4Gpio, 0, true, nil, make(chan int)}
	return
}

// ArmStart sets up the interrupt handler for the start GPIO line
func ArmStart() (start GpioTime, err error) {
	start = GpioTime{startChip, startGpio, 0, true, nil, make(chan int)}
	start.Line, err = gpiod.RequestLine(start.Chip, start.Lane, gpiod.AsInput,
		gpiod.WithEventHandler(start.gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	return
}

// ArmLanes sets up the interrupt handler for the all the lane GPIO lines
func ArmLanes() (lanes [4]GpioTime, err error) {
	lanes = createLanes()
	lanes[0].Line, err = gpiod.RequestLine(lanes[0].Chip, lanes[0].Lane, gpiod.AsInput,
		gpiod.WithEventHandler(lanes[0].gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	if err != nil {
		log.Fatal(err)
	}
	lanes[1].Line, err = gpiod.RequestLine(lanes[1].Chip, lanes[1].Lane, gpiod.AsInput,
		gpiod.WithEventHandler(lanes[1].gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	if err != nil {
		log.Fatal(err)
	}
	lanes[2].Line, err = gpiod.RequestLine(lanes[2].Chip, lanes[2].Lane, gpiod.AsInput,
		gpiod.WithEventHandler(lanes[2].gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	if err != nil {
		log.Fatal(err)
	}
	lanes[3].Line, err = gpiod.RequestLine(lanes[3].Chip, lanes[3].Lane, gpiod.AsInput,
		gpiod.WithEventHandler(lanes[3].gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	if err != nil {
		log.Fatal(err)
	}
	return
}

// WaitForStart waits until the start GPIO triggers
func WaitForStart(start *GpioTime) {
	pending := start.Pending
	for pending {
		select {
		case <-start.Channel:
			pending = false
		}
	}
	start.Close()
}

// deltaTimes calculates the difference betwene two timestamps
func deltaTimes(start GpioTime, end GpioTime) float64 {
	if end.Pending || start.Pending {
		return 0.0
	}
	return end.Time.Seconds() - start.Time.Seconds()
}

// WaitForLanes waits until all 4 lanes have triggered and returns
// the time difference for each lane
func WaitForLanes(lanes [4]GpioTime) {
	done := [4]bool{lanes[0].Pending, lanes[1].Pending, lanes[2].Pending, lanes[3].Pending}
	for !done[0] && !done[1] && !done[2] && !done[3] {
		select {
		case <-lanes[0].Channel:
			lanes[0].Close()
			done[0] = true
			log.Println("Lane 0 done.")
		case <-lanes[1].Channel:
			lanes[1].Close()
			done[1] = true
			log.Println("Lane 1 done.")
		case <-lanes[2].Channel:
			lanes[2].Close()
			done[2] = true
			log.Println("Lane 2 done.")
		case <-lanes[3].Channel:
			lanes[3].Close()
			done[3] = true
			log.Println("Lane 3 done.")
		case <-time.After(20 * time.Second):
			log.Println("Lanes have timed out.")
			for i, _ := range done {
				done[i] = true
			}
		}
	}
	for _, g := range lanes {
		g.Close()
	}
	return
}

// GetTimes returns the difference between a set of lanes and start time
// GpioTime structures
func GetTimes(start GpioTime, lanes [4]GpioTime) (times [4]float64) {
	for i, _ := range times {
		times[i] = deltaTimes(start, lanes[i])
	}
	return
}
