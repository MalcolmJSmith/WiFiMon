package main

import (
	"fmt"
	"os"
	"time"
	"runtime"
)

type AboutInfo struct {
	ProgVersion string
	GoVersion string
}

func MessageLoop() {
	var (
		idle   int
		s      WiFiStatus
		a AboutInfo
	)
	timer := time.Tick(time.Second)
	a.ProgVersion = "1.0"
	a.GoVersion = runtime.Version()

	for {
		select {
		case _ = <-timer:
			if idle < 5 {
				// Call QueryMonitor once a second while there is an active client
				idle++
				select {
				case chMon <- s:
				default:
				}
			}
		case s = <-chMon: 		
				select {
				case _ = <- chExit:
					StopMonitor()					
					fmt.Println("Shutdown request received")
					os.Exit(0)
				default:
				}
		case chStatus <- s:
			// A request has been served to a client
			idle = 0
		case chAbout <- a:
			
		}
	}
}
