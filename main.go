package main

import (
	_ "embed"
	"github.com/Nescient/gpio-timer/derbynet"
	"github.com/Nescient/gpio-timer/gpio"
	"log"
	"os"
	"runtime/debug"
	// "sync"
	"time"
)

//go:generate sh -c "printf %s $(git rev-parse HEAD) > .commit_id"
//go:embed .commit_id
var gitrev string

func main() {

	whoAmI := os.Args[0]
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				whoAmI += setting.Value
			} else if setting.Key == "vcs.time" {
				whoAmI += setting.Value
			} else if setting.Key == "vcs.modified" {
				whoAmI += "-DIRTY"
			}
		}
	}

	log.Println(whoAmI)
	log.Println(gitrev)

	derbynet.GetCookie()
	derbynet.Hello()
	derbynet.Identified(gitrev)

	// timer heartbeats
	isQuitting := false
	go func() {
		for !isQuitting {
			derbynet.Heartbeat()
			time.Sleep(time.Second * 5)
		}
		derbynet.Terminate()
	}()

	// main race loop
	for isQuitting == false {
		log.Println("Waiting for heat...")
		if derbynet.WaitForHeat() {
			l, err := gpio.ArmStart()
			if err != nil {
				// log.Fatal(err)
				log.Println(err)
			}
			log.Println("Waiting for start gate...")
			gpio.WaitForStart()
			l.Close()
			derbynet.Started()
			ll, err := gpio.ArmLanes()
			if err != nil {
				// log.Fatal(err)
				log.Println(err)
			}
			log.Println("Waiting for lanes...")
			laneTimes := gpio.WaitForLanes()
			ll.Close()
			derbynet.Finished(laneTimes[0], laneTimes[1], laneTimes[2], laneTimes[3])
			log.Println("Times %f %f %f %f", laneTimes[0], laneTimes[1], laneTimes[2], laneTimes[3])
		}
	}

	x, y := gpio.GetGateTime()
	log.Printf("gate time is %v, %d\n", x, y)

	// var wg = &sync.WaitGroup{}
	// l, err := gpio.Arm()
	// if err != nil {
	// log.Fatal(err)
	// }
	// log.Println(l)

	// wg.Add(1)
	// wg.Wait()
	// gpio.WaitForStart()

	isQuitting = true
	log.Println("Terminating...")
	derbynet.Terminate()
	time.Sleep(time.Second * 2)
}
