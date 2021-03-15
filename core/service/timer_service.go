package service

import (
	"fmt"
	"time"
)

type ScheduledFunction func() error

var done = make(chan bool, 1)
var ticker *time.Ticker

func InitializeTimer(ticksInSeconds int, scheduledFunction ScheduledFunction) {
	timerTick := time.Duration(ticksInSeconds) * time.Second

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