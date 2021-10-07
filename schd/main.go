package main

import (
	"log"
	"strings"
	"time"
)

func main() {
	// job := func() {
	// 	t := time.Now()
	// 	fmt.Println("Time's up! @", t.UTC())
	// }
	// // Run every 2 seconds but not now.
	// //scheduler.Every(2).Seconds().NotImmediately().Run(job)

	// // Run now and every X.
	// scheduler.Every(1).Minutes().Run(job)
	// //scheduler.Every().Day().Run(job)
	// //scheduler.Every().Monday().At("08:30").Run(job)

	// // Keep the program from not exiting.
	// runtime.Goexit()

	// wg := sync.WaitGroup{}

	// ch := make(chan string, 10)
	// for i := 0; i < 10; i++ {
	// 	wg.Add(1)
	// 	go sum(&wg, ch)
	// }
	// wg.Wait()
	// close(ch)

	// for i := range ch {
	// 	println(i)
	// }

	layout := "2006-01-02 15:04:05"
	localLoc, err := time.LoadLocation("America/Los_Angeles") //time.LoadLocation("Europe/Paris")
	if err != nil {
		log.Fatal("Failed to load location")
	}
	bt, _ := time.Parse(layout, "2021-03-01 00:00:00")
	localDateTime := bt.In(localLoc)

	days := time.Now().Sub(localDateTime).Hours() / 24
	println(int(days))
	if days > 7 {
		baseTime := localDateTime.AddDate(0, 0, 5).Format(layout)
		println(baseTime)
	}
	sql := "asdfasdfasdfa"
	ip := []string{"100.141", "100.142"}
	sql += "('" + strings.Join(ip, "','") + "')"
	println(sql)
}
