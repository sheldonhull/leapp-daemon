package timer

import (
	"time"
)

type timer struct {
	done   chan bool
	ticker *time.Ticker
}

type ScheduledFunction func()

func NewTimer(intervalInSeconds int, scheduledFunction ScheduledFunction) *timer {
	tickDuration := time.Duration(intervalInSeconds) * time.Second

	timer := timer{
		done:   make(chan bool),
		ticker: time.NewTicker(tickDuration),
	}

	go func() {
		for {
			select {
			case <-timer.done:
				return
			case <-timer.ticker.C:
				scheduledFunction()
			}
		}
	}()

	return &timer
}

func (timer *timer) Close() {
	timer.ticker.Stop()
	timer.done <- true
}
