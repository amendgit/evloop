package evloop

import (
	"log"
	"time"
)

type tEventPump struct {
	evloop             *EventLoop
	nextDelayedRunTime time.Time
	hasPendingEvent    bool
}

func newEventPump(evloop *EventLoop) *tEventPump {
	var pump = new(tEventPump)
	pump.evloop = evloop
	pump.nextDelayedRunTime = time.Now()
	return pump
}

func (pump *tEventPump) scheduleEvent() {
	pump.hasPendingEvent = true
}

func (pump *tEventPump) scheduleDelayedEvent(delayedRunTime time.Time) {
	var now = time.Now()

	if pump.nextDelayedRunTime.Before(now) {
		pump.nextDelayedRunTime = delayedRunTime
		return
	}

	if pump.nextDelayedRunTime.After(delayedRunTime) && delayedRunTime.After(now) {
		pump.nextDelayedRunTime = delayedRunTime
	}
}

func (pump *tEventPump) run() {
	var continueCount = 0

	for {
		var moreWorkIsPlausible, more = false, false
		pump.hasPendingEvent = false

		more = pump.evloop.processEvent()
		moreWorkIsPlausible = moreWorkIsPlausible || more
		if pump.evloop.shouldQuit {
			break
		}

		more, nextDelayedRunTime := pump.evloop.processDelayedEvent()
		moreWorkIsPlausible = moreWorkIsPlausible || more
		if pump.evloop.shouldQuit {
			break
		}

		pump.scheduleDelayedEvent(nextDelayedRunTime)
		moreWorkIsPlausible = moreWorkIsPlausible || pump.hasPendingEvent
		if moreWorkIsPlausible {
			continueCount++
			if continueCount >= 1000 {
				log.Printf("to much pending events.")
				time.Sleep(time.Second)
			}
			continue
		}
		continueCount = 0

		var wait = pump.nextDelayedRunTime.Sub(time.Now())
		if wait > 0 {
			time.Sleep(wait)
		}
	}
}
