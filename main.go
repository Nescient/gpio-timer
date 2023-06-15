package main

import (
	"fmt"
	"github.com/Nescient/gpio-timer/derbynet"
	"sort"
	"strings"
	"time"
)

func main() {
	s := []int{345, 78, 123, 10, 76, 2, 567, 5}
	sort.Ints(s)
	fmt.Println("Sorted slice: ", s)

	// Finding the index
	fmt.Println("Index value: ", strings.Index("GeeksforGeeks", "ks"))

	// Finding the time
	fmt.Println("Time: ", time.Now().Unix())

	derbynet.GetCookie()
	derbynet.Hello()
	derbynet.Identified()

	for i := 0; i < 10; i++ {
		derbynet.Heartbeat()
		time.Sleep(time.Second * 5)
	}
}
