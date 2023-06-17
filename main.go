package main

import (
	_ "embed"
	"github.com/Nescient/gpio-timer/derbynet"
	"github.com/Nescient/gpio-timer/gpio"
	"log"
	"os"
	"runtime/debug"
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

	for i := 0; i < 2; i++ {
		derbynet.Heartbeat()
		derbynet.Started()
		time.Sleep(time.Second * 5)
		derbynet.Finished(12.34567890, 15.678901234, 0, 9.99999)
	}

	x, y := gpio.GetGateTime()
	log.Printf("gate time is %v, %d\n", x, y)

	l, err := gpio.Arm()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(l)
}
