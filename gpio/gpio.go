// gpio is a wrapper package using warthog618's gpiod to watch specific GPIO devices for
// changes.  this is the main functionality of gpio-timer
package gpio

import (
	"github.com/loov/hrtime"
	"github.com/warthog618/gpiod"
	"log"
	"sync"
	"time"
)

// pin mapping, see https://hub.libre.computer/t/libre-computer-wiring-tool/40
// Pin     Chip    Line    sysfs   Name    Pad     Ref     Desc
// 1       3.3V    3.3V    3.3V    3.3V    3.3V    VCC_IO  3.3V
// 2       5V      5V      5V      5V      5V      VCC_SYS 5V
// 3       2       25      89      GPIO2_D1        R17     I2C0_SDA        I2C0_SDA/FEPHY_LED_DATA_M1
// 4       5V      5V      5V      5V      5V      VCC_SYS 5V
// 5       2       24      88      GPIO2_D0        P17     I2C0_SCL        I2C0_SCL/FEPHY_LED_LINK_M1
// 6       GND     GND     GND     GND     GND     GND     GND
// 7       1       28      60      GPIO1_D4        Y17     GPIO1_D1        CLKOUT
// 8       3       4       100     GPIO3_A4        E2      UART1_TX        TSP_D0/CIF_D0/SDMMC0EXT_D0/UART1_TX/USB3PHY_DEBUG4_u
// 9       GND     GND     GND     GND     GND     GND     GND
// 10      3       6       102     GPIO3_A6        F2      UART1_RX        TSP_D2/CIF_D2/SDMMC0EXT_D2/UART1_RX/USB3PHY_DEBUG6_u
// 11      2       20      84      GPIO2_C4        V18     GPIO2_C4_U/I2S1_SDO1    I2S1_SDIO1/PDM_SDI1_M0/CARD_RST_M1_u
// 12      2       6       70      GPIO2_A6        M19     GPIO2_A6_U/PWM2 PWM2_u
// 13      2       21      85      GPIO2_C5        V17     GPIO2_C5_U/I2S1_SDO2    I2S1_SDIO2/PDM_SDI2_M0/CARD_DET_M1_u
// 14      GND     GND     GND     GND     GND     GND     GND
// 15      2       22      86      GPIO2_C6        V16     GPIO2_C6_U/I2S1_SDO3    I2S1_SDIO3/PDM_SDI3_M0/CARD_IO_M1_u
// 16      3       7       103     GPIO3_A7        F1      UART1_CTSN      TSP_D3/CIF_D3/SDMMC0EXT_D3/UART1_CTSN/USB3PHY_DEBUG7_u
// 17      GND     GND     GND     GND     GND     VCC_IO  GND
// 18      3       5       101     GPIO3_A5        D1      UART1_RTSN      TSP_D1/CIF_D1/SDMMC0EXT_D1/UART1_RTSN/USB3PHY_DEBUG5_u
// 19      3       1       97      GPIO3_A1        D2      GPIO3_A1_U/SPI_TXD      TSP_FAIL/CIF_HREF/SDMMC0EXT_DET/SPI_TXD_M2/USB3PHY_DEBUG2/I2S2_SDO_M1_u
// 20      GND     GND     GND     GND     GND     GND     GND
// 21      3       2       98      GPIO3_A2        E1      GPIO3_A2_D/SPI_RXD      TSP_CLK/CIF_CLKIN/SDMMC0EXT_CLK/SPI_RXD_M2/USB3PHY_DEBUG3/I2S2_SDI_M1_d
// 22      0       2       2       GPIO0_A2        R3      GPIOA2/CLKOUT/SPDIF_TX_M2       CLKOUT_GMAC_M0/SPDIF_TX_M2_d
// 23      3       0       96      GPIO3_A0        E3      GPIO3_A0_U/SPI_CLK      TSP_VALID/CIF_VSYNC/SDMMC0EXT_CMD/SPI_CLK_M2/USB3PHY_DEBUG1/I2S2_SCLK_M1_u
// 24      3       8       104     GPIO3_B0        F3      GPIO3_B0_D/SPI_CSN0     TSP_D4/CIF_D4/SPI_CSN0_M2/I2S2_LRCK_TX_M1/USB3PHY_DEBUG8/I2S2_LRCK_RX_M1_d
// 25      GND     GND     GND     GND     GND     GND     GND
// 26      2       12      76      GPIO2_B4        T16     GPIO2_B4/SPI_CSN1_M0/FLASH_VOL_SEL      SPI_CSN1_M0/FLASH_VOL_SEL_u
// 27      2       4       68      GPIO2_A4        N19     I2C1_SDA_PMIC   I2C1_SDA
// 28      2       5       69      GPIO2_A5        N20     I2C1_SCL_PMIC   I2C1_SCL
// 29      2       19      83      GPIO2_C3        U16     GPIO2_C3_U/I2S1_SDI     I2S1_SDI/PDM_SDI0_M0/CARD_CLK_M1_u
// 30      GND     GND     GND     GND     GND     GND     GND
// 31      2       23      87      GPIO2_C7        N17     GPIO2_C7_U/I2S1_SDO     I2S1_SDO/PWDM_FSYNC_M0_u
// 32      0       0       0       GPIO0_A0        L3      GPIO0_A0/CLKOUT CLKOUT_WIFI_M0_d
// 33      2       16      80      GPIO2_C0/GPIO2_C1*      V15/P18 GPIO2_C0_U/I2S1_LRCK_RX / GPIO2_C1_U/I2S1_LRCK_TX       I2S1_LRCK_RX/TSP_D5_M1/CIF_D5_M1_u / I2S1_LRCK_TX/SPDIF_TX_M1/TSP_D6_M1/CIF_D6_M1_u
// 34      GND     GND     GND     GND     GND     GND     GND
// 35      2       18      82      GPIO2_C2        R18     GPIO2_C2_D/I2S1_SCLK    I2S1_SCLK/PDM_CLK_M0/TSP_D7_M1/CIF_D7_M1_d
// 36      2       0       64      GPIO2_A0        P19     DEBUG_TX        UART2_TX_M1/POWERSTATE0_d
// 37      2       15      79      GPIO2_B7        N18     GPIO2_B7_D/I2S1_MCLK    I2S1_MCLK/TSP_SYNC_M1/CIF_CLKOUT_M1_d
// 38      2       1       65      GPIO2_A1        P20     DEBUG_RX        UART2_RX_M1/POWERSTATE1_u
// 39      GND     GND     GND     GND     GND     GND     GND
// 40      0       27      27      GPIO0_D3        V9      SPDIF_TX_M0     SPDIF_TX_M0_d

var startChip = "gpiochip3"
var laneChip = "gpiochip2"
var startGpio = 4
var lane1Gpio = 24
var lane2Gpio = 25
var lane3Gpio = 20
var lane4Gpio = 21

type Timestamp struct {
	Time     time.Duration
	Count    hrtime.Count
	GpioTime time.Duration
}

// var startTime time.Duration
// var startCount hrtime.Count
var startTime Timestamp

var waitStart sync.WaitGroup
var waitLanes sync.WaitGroup

var laneTimes = [4]Timestamp{}

// type DerbyTimer struct {
//    startChip = "gpiochip1"
// laneChip = "gpiochip2"
// startGpio = 28
// lane1Gpio = 24
// lane2Gpio = 25
// lane3Gpio = 20
// lane4Gpio = 21
// waitGroup sync.WaitGroup
// StartTime time.Duration
// StartCount hrtime.Count
// }

// func (this DerbyTimer) Arm() err error {
// this.waitGroup.Add(1)
// }

// init will
func init() {

}

func clearLanes() {
	laneTimes = [4]Timestamp{Timestamp{0, 0, 0}, Timestamp{0, 0, 0}, Timestamp{0, 0, 0}, Timestamp{0, 0, 0}}
}

// setStartTime sets the time that the gate started
func setStartTime(evt gpiod.LineEvent) {
	gpioNum := evt.Offset // an int
	time := evt.Timestamp // time.Duration
	startTime.Time = hrtime.Now()
	startTime.Count = hrtime.TSC()
	startTime.GpioTime = evt.Timestamp
	log.Printf("got event %d, expecting %d\n", gpioNum, startGpio)
	log.Printf("got gate start at %v, %v, %d\n", time, startTime.Time, startTime.Count)
	waitStart.Done()
}

// setLaneTime sets the time that a given lane completes
func setLaneTime(evt gpiod.LineEvent) {
	log.Printf("got lane event %d\n", evt.Offset)
	switch gpioNum := evt.Offset; gpioNum {
	case lane1Gpio:
		if laneTimes[0].Count == 0 {
			laneTimes[0] = Timestamp{hrtime.Now(), hrtime.TSC(), evt.Timestamp}
			waitLanes.Done()
		}
	case lane2Gpio:
		if laneTimes[1].Count == 0 {
			laneTimes[1] = Timestamp{hrtime.Now(), hrtime.TSC(), evt.Timestamp}
			waitLanes.Done()
		}
	case lane3Gpio:
		if laneTimes[2].Count == 0 {
			laneTimes[2] = Timestamp{hrtime.Now(), hrtime.TSC(), evt.Timestamp}
			waitLanes.Done()
		}
	case lane4Gpio:
		if laneTimes[3].Count == 0 {
			laneTimes[3] = Timestamp{hrtime.Now(), hrtime.TSC(), evt.Timestamp}
			waitLanes.Done()
		}
	default:
		log.Printf("unknown lane event %d\n", gpioNum)
	}
	// log.Printf("got event %d, expecting %d\n", gpioNum, startGpio)
	// log.Printf("got gate start at %v, %v, %d\n", time, startTime, startCount)
	// waitLanes.Done()
}

func ArmStart() (*gpiod.Line, error) {
	clearLanes()
	waitStart.Add(1)
	// gpiod.WithBothEdges and then we wont care really ?
	return gpiod.RequestLine(startChip, startGpio, gpiod.AsInput,
		gpiod.WithEventHandler(setStartTime), gpiod.LineEdgeFalling)
}

func ArmLanes() (*gpiod.Lines, error) {
	clearLanes()
	waitLanes.Add(4)
	// gpiod.WithBothEdges and then we wont care really ?
	return gpiod.RequestLines(laneChip, []int{lane1Gpio, lane2Gpio, lane3Gpio, lane4Gpio},
		gpiod.AsInput, gpiod.WithEventHandler(setLaneTime), gpiod.LineEdgeFalling)
}

func WaitForStart() {
	waitStart.Wait()
}

func deltaTimes(start Timestamp, end Timestamp) float64 {
	if end.Time < start.Time {
		return 0.0
	}
	// delta1 := end.Time.Sub(start.Time).Seconds()
	delta1 := end.Time.Seconds() - start.Time.Seconds()
	log.Printf("delta TSC is %d, %d\n", start.Count, end.Count)
	delta2 := (end.Count - start.Count).ApproxDuration().Seconds()
	delta3 := end.GpioTime.Seconds() - start.GpioTime.Seconds()
	log.Printf("delta times %f, %f, %f \n", delta1, delta2, delta3)
	return delta1
}

func WaitForLanes() (times [4]float64) {
	waitLanes.Wait()
	for i, _ := range laneTimes {
		times[i] = deltaTimes(startTime, laneTimes[i])
	}
	return
	// times[0] = deltaTimes(startTime, laneTimes[0].Time)
	// times[0] = deltaTimes(startTime, laneTimes[0].Time)
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
