// Copyright 2015 By Jash. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// mail: shijian0912@163.com

package evloop

import (
	"container/list"
	"sync"
	"time"
)

type EventLoop struct {
	shouldQuit             bool
	pendingEventQueue      *list.List
	pendingEventQueueMutex sync.Mutex
	delayedEventQueue      *Pqueue
	pump                   *tEventPump
	recentTime             time.Time
}

func NewEventLoop() *EventLoop {
	var el = new(EventLoop)
	el.pendingEventQueue = list.New()
	el.delayedEventQueue = NewPqueue()
	el.pump = newEventPump(el)
	return el
}

func (el *EventLoop) PostEvent(event *Event) {
	if event.DelayedRunTime.After(time.Now()) {
		el.addToDelayedEventQueue(event)
	} else {
		el.addToPendingEventQueue(event)
	}
}

func (el *EventLoop) PostFunc(function func()) {
	var event = NewEvent(function, 0)
	el.PostEvent(event)
}

func (el *EventLoop) PostDelayedFunc(function func(), delayed time.Duration) {
	var event = NewEvent(function, delayed)
	el.PostEvent(event)
}

func (el *EventLoop) RepeatFunc(function func(stop *bool), delayed time.Duration) {
	var repeat *Event
	repeat = NewEvent(func() {
		repeat.DelayedRunTime = time.Now().Add(delayed)
		var stop bool
		function(&stop)
		if stop {
			return
		}
		el.PostEvent(repeat)
	}, 0)

	el.PostEvent(repeat)
}

func (el *EventLoop) ShouldQuit() {
	el.shouldQuit = true
}

func (el *EventLoop) addToPendingEventQueue(event *Event) {
	var wasEmpty = (el.pendingEventQueue.Len() == 0)
	el.pendingEventQueue.PushBack(event)
	el.scheduleEvent(wasEmpty)
}

func (el *EventLoop) addToDelayedEventQueue(event *Event) {
	el.delayedEventQueue.Push(event)
	var nextDelayedRunTime = el.delayedEventQueue.Top().(*Event).DelayedRunTime
	el.scheduleDelayedEvent(nextDelayedRunTime)
}

func (el *EventLoop) scheduleEvent(wasEmpty bool) {
	if wasEmpty {
		el.pump.scheduleEvent()
	}
}

func (el *EventLoop) scheduleDelayedEvent(dealyedRunTime time.Time) {
	el.pump.scheduleDelayedEvent(dealyedRunTime)
}

func (el *EventLoop) processEvent() (handled bool) {
	var now = time.Now()

	for el.pendingEventQueue.Len() != 0 {
		el.pendingEventQueueMutex.Lock()
		var elem = el.pendingEventQueue.Front()
		var event = elem.Value.(*Event)
		el.pendingEventQueue.Remove(elem)
		el.pendingEventQueueMutex.Unlock()

		if !event.DelayedRunTime.After(now) {
			event.Function()
			return el.pendingEventQueue.Len() != 0
		} else {
			el.delayedEventQueue.Push(event)
			if el.delayedEventQueue.Top() == Interface(event) {
				el.pump.scheduleDelayedEvent(el.delayedEventQueue.Top().(*Event).DelayedRunTime)
			}
		}
	}

	return false
}

func (el *EventLoop) processDelayedEvent() (handled bool, nextDelayedEventTime time.Time) {
	if el.delayedEventQueue.Empty() {
		el.recentTime = time.Now()
		return false, el.recentTime
	}

	// When we "fall behind", there will be a lot of tasks in the delayed event
	// queue that are ready to run. To increase efficiency when we fall behind,
	// we will only call Time.Now() intermittently, and then process all tasks
	// that are ready to run before calling it again. As a result, the more we
	// fall behind (and have a lot of ready-to-run delayed events), the more
	// efficiency we'll be at handling the events.
	var nextRunTime = el.delayedEventQueue.Top().(*Event).DelayedRunTime
	if nextRunTime.After(el.recentTime) {
		el.recentTime = time.Now()
		if nextRunTime.After(el.recentTime) {
			return false, nextRunTime
		}
	}

	var event = el.delayedEventQueue.Top().(*Event)
	el.delayedEventQueue.Pop()

	event.Function()

	if !el.delayedEventQueue.Empty() {
		nextDelayedEventTime = el.delayedEventQueue.Top().(*Event).DelayedRunTime
	}

	return true, nextDelayedEventTime
}

func (el *EventLoop) Run() {
	el.pump.run()
}
