package service

import (
	"fmt"
	"time"
)

type ScheduledFunction func() error

var done = make(chan bool, 1)
var ticker *time.Ticker

func InitiliazeTimer(ticksInSeconds int, scheduledFunction ScheduledFunction) {
	timerTick := time.Duration(ticksInSeconds) * time.Second // 1 sec

	if ticker == nil {
		ticker = time.NewTicker(timerTick)
	}

	if done == nil {
		done = make(chan bool)
	}

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				err := scheduledFunction()
				if err != nil {
					panic(err)
				}
				fmt.Println("Tick at", t)
			}
		}
	}()
}

func CloseTimer() {
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}