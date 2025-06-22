// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/warthog618/gpiod"
	"log"
	"sync/atomic"
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
var lane2Gpio = 24
var lane3Gpio = 25
var lane4Gpio = 20
var lane1Gpio = 21

// GpioTime is a structure that represents a single GPIO event
type GpioTime struct {
	Chip    string
	Lane    int
	Time    time.Duration
	Pending atomic.Bool
	Line    *gpiod.Line
	Channel chan int
}

// New initializes the structure to default values
func (this *GpioTime) New(chip string, offset int) {
	this.Close()
	this.Chip = chip
	this.Lane = offset
	this.Time = 0
	this.Pending.Store(true)
	this.Line = nil
	this.Channel = make(chan int)
}

// Arm will register the GPIO for a falling edge event
func (this *GpioTime) Arm(l bool) (err error) {
	this.Pending.Store(true)
	if l {
		log.Printf("doing rising edge\n")
           this.Line, err = gpiod.RequestLine(this.Chip, this.Lane, gpiod.AsInput,
	      gpiod.WithEventHandler(this.gpioHandler), gpiod.LineEdgeRising, gpiod.WithPullDown,
	      gpiod.WithDebounce(1*time.Millisecond))
      }	else {
	      log.Printf("doing falling edge\n")
	   this.Line, err = gpiod.RequestLine(this.Chip, this.Lane, gpiod.AsInput,
		gpiod.WithEventHandler(this.gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp,
		gpiod.WithDebounce(1*time.Millisecond))
	}
	return
}

// gpioHandler handles a GPIO event for a given GpioTime struct
func (this *GpioTime) gpioHandler(evt gpiod.LineEvent) {
	if evt.Offset == this.Lane {
		// if pending, swap and set time
		if this.Pending.CompareAndSwap(true, false) {
			log.Printf("Received GPIO event %d at %v\n", evt.Offset, evt.Timestamp)
			this.Time = evt.Timestamp
			// need to non-blocking send this
			select {
			case this.Channel <- 1:
				// message sent
			default:
				// message dropped
			}
		} else {
			log.Printf("Received GPIO event %d twice\n", evt.Offset)
		}
	} else {
		log.Printf("Received unknown GPIO event %d\n", evt.Offset)
	}
	log.Println("..end handler..")
}

// WaitForever will wait until the handler is called
func (this *GpioTime) WaitForever() {
	for this.Pending.Load() {
		select {
		case <-this.Channel:
			log.Printf("GPIO %d complete.", this.Lane)
		case <-time.After(1 * time.Second):
		}
	}
	log.Println("..end wait..")
}

// WaitFor will wait until the handler is called or a set
// amount of time expires
func (this *GpioTime) WaitFor(timeout time.Duration) {
	if timeout > 0 && this.Pending.Load() {
		select {
		case <-this.Channel:
			log.Printf("GPIO %d complete.", this.Lane)
		case t := <-time.After(timeout):
			log.Printf("GPIO %d timeout at %v.\n", this.Lane, t)
		}
	}
}

// Close will close any open GPIO lanes for the GpioTime struct
// as well as the channel
func (this *GpioTime) Close() {
	if this.Line != nil {
		this.Line.Close()
		this.Line = nil
	}
	// close(this.Channel) // not safe to do multiple times
	log.Println("..end close..")
}

// createLanes initializes an array of GpioTime structures
// to represent the set of Gpio lanes
func createLanes() (lanes [4]*GpioTime) {
	for i, _ := range lanes {
		lanes[i] = new(GpioTime)
	}
	lanes[0].New(laneChip, lane1Gpio)
	lanes[1].New(laneChip, lane2Gpio)
	lanes[2].New(laneChip, lane3Gpio)
	lanes[3].New(laneChip, lane4Gpio)
	return
}

// ArmStart sets up the interrupt handler for the start GPIO line
func ArmStart() (start *GpioTime, err error) {
	start = new(GpioTime)
	start.New(startChip, startGpio)
	err = start.Arm(false)
	return
}

// ArmLanes sets up the interrupt handler for the all the lane GPIO lines
func ArmLanes() (lanes [4]*GpioTime, err error) {
	lanes = createLanes()
	for i, _ := range lanes {
		err = lanes[i].Arm(false)
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

// WaitForStart waits until the start GPIO triggers
func WaitForStart(start *GpioTime) {
	start.WaitForever()
	start.Close()
}

// deltaTimes calculates the difference betwene two timestamps
func deltaTimes(start *GpioTime, end *GpioTime) float64 {
	if end.Pending.Load() || start.Pending.Load() {
		return 0.0
	}
	return end.Time.Seconds() - start.Time.Seconds()
}

// WaitForLanes waits until all 4 lanes have triggered and returns
// the time difference for each lane
func WaitForLanes(lanes [4]*GpioTime) {
	doneAt := time.Now().Add(15 * time.Second)
	for i, _ := range lanes {
		lanes[i].WaitFor(doneAt.Sub(time.Now()))
		lanes[i].Close()
	}
}

// GetTimes returns the difference between a set of lanes and start time
// GpioTime structures
func GetTimes(start *GpioTime, lanes [4]*GpioTime) (times [4]float64) {
	for i, _ := range times {
		times[i] = deltaTimes(start, lanes[i])
	}
	return
}
