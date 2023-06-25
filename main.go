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

	var client derbynet.DerbyNet
	client.Initialize()
	log.Println("Getting cookie...")
	client.GetCookie()
	log.Println("Saying hello...")
	client.Hello()
	log.Println("Indentifying...")
	client.Identified(gitrev)

	// timer heartbeats
	log.Println("Establishing heartbeats...")
	isQuitting := false
	go func() {
		for !isQuitting {
			client.Heartbeat()
			time.Sleep(time.Second * 5)
		}
		client.Terminate()
	}()

	// main race loop
	log.Println("Starting main race loop...")
	for isQuitting == false {
		log.Println("Waiting for heat...")
		if client.WaitForHeat() {
			start, err := gpio.ArmStart()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Waiting for start gate...")
			gpio.WaitForStart(start)
			client.Started()
			lanes, err := gpio.ArmLanes()
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Waiting for lanes...")
			gpio.WaitForLanes(lanes)
			laneTimes := gpio.GetTimes(start, lanes)
			client.Finished(laneTimes[0], laneTimes[1], laneTimes[2], laneTimes[3])
			log.Printf("Times %f %f %f %f\n", laneTimes[0], laneTimes[1], laneTimes[2], laneTimes[3])
		}
	}

	isQuitting = true
	log.Println("Terminating...")
	client.Terminate()
	time.Sleep(time.Second * 2)
}
