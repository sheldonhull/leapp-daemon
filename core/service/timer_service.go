package service

import (
	"leapp_daemon/shared/logging"
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
			case <-ticker.C:
				err := scheduledFunction()
				if err != nil {
					logging.Entry().Error(err)
					panic(err)
				}
			}
		}
	}()
}

func CloseTimer() {
	ticker.Stop()
	done <- true
	logging.Info("Ticker stopped")
}