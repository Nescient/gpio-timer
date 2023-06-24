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
func (this *GpioTime) New(chip string, offset int) {
	this.Close()
	this.Chip = chip
	this.Lane = offset
	this.Time = 0
	this.Pending = true
	this.Line = nil
	this.Channel = make(chan int)
}

// Arm will register the GPIO for a falling edge event
func (this *GpioTime) Arm() (err error) {
	this.Pending = true
	this.Line, err = gpiod.RequestLine(this.Chip, this.Lane, gpiod.AsInput,
		gpiod.WithEventHandler(this.gpioHandler), gpiod.LineEdgeFalling, gpiod.WithPullUp)
	return
}

// gpioHandler handles a GPIO event for a given GpioTime struct
func (this *GpioTime) gpioHandler(evt gpiod.LineEvent) {
	if evt.Offset == this.Lane {
		this.Pending = false
		this.Time = evt.Timestamp
		log.Printf("Received GPIO event at %v\n", this.Time)
		this.Channel <- 1
	} else {
		log.Printf("Received unknown GPIO event %d\n", evt.Offset)
	}
}

// Wait will wait until the handler is called or a set
// amount of time expires
func (this *GpioTime) Wait(seconds int) {
	pending := start.Pending
	if seconds > 0 {
		for pending {
			select {
			case <-start.Channel:
				pending = false
			case <-time.After(seconds * time.Second):
				return
			}
		}
	} else {
		for pending {
			select {
			case <-start.Channel:
				pending = false
			}
		}
	}
}

// Close will close any open GPIO lanes for the GpioTime struct
// as well as the channel
func (this *GpioTime) Close() {
	if this.Line != nil {
		this.Line.Close()
	}
	// close(this.Channel) // not safe to do multiple times
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
	err = start.Arm()
	return
}

// ArmLanes sets up the interrupt handler for the all the lane GPIO lines
func ArmLanes() (lanes [4]*GpioTime, err error) {
	lanes = createLanes()
	for i, _ := range lanes {
		err = lanes[i].Arm()
		if err != nil {
			log.Fatal(err)
		}
	}
	return
}

// WaitForStart waits until the start GPIO triggers
func WaitForStart(start *GpioTime) {
	start.Wait(0) // wait forever
	start.Close()
	log.Printf("wait for statr %v\n", start.Time)
}

// deltaTimes calculates the difference betwene two timestamps
func deltaTimes(start *GpioTime, end *GpioTime) float64 {
	if end.Pending || start.Pending {
		return 0.0
	}
	return end.Time.Seconds() - start.Time.Seconds()
}

// WaitForLanes waits until all 4 lanes have triggered and returns
// the time difference for each lane
func WaitForLanes(lanes [4]*GpioTime) {
	done := [4]bool{
		!lanes[0].Pending,
		!lanes[1].Pending,
		!lanes[2].Pending,
		!lanes[3].Pending,
	}
	for !done[0] || !done[1] || !done[2] || !done[3] {
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
func GetTimes(start *GpioTime, lanes [4]*GpioTime) (times [4]float64) {
	for i, _ := range times {
		times[i] = deltaTimes(start, lanes[i])
	}
	return
}
