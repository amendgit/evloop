// Copyright 2015 By Jash. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// mail: shijian0912@163.com

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
		var shouldContinue, more = false, false
		pump.hasPendingEvent = false

		more = pump.evloop.processTask()
		shouldContinue = shouldContinue || more
		if pump.evloop.shouldQuit {
			break
		}

		more, nextDelayedRunTime := pump.evloop.processDelayedTask()
		shouldContinue = shouldContinue || more
		if pump.evloop.shouldQuit {
			break
		}

		shouldContinue = shouldContinue || pump.hasPendingEvent
		if shouldContinue {
			continueCount++
			if continueCount >= 1000 {
				log.Printf("to much pending events.")
				time.Sleep(time.Second)
			}
			continue
		}
		continueCount = 0

		pump.scheduleDelayedEvent(nextDelayedRunTime)
		var wait = pump.nextDelayedRunTime.Sub(time.Now())
		if wait > 0 {
			time.Sleep(wait)
		}
	}
}
