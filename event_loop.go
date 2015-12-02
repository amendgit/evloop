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
	shouldQuit            bool
	pendingTaskQueue      *list.List
	pendingTaskQueueMutex sync.Mutex
	delayedTaskQueue      *Pqueue
	pump                  *tEventPump
	recentTime            time.Time
}

func NewEventLoop() *EventLoop {
	var el = new(EventLoop)
	el.pendingTaskQueue = list.New()
	el.delayedTaskQueue = NewPqueue()
	el.pump = newEventPump(el)
	return el
}

func (el *EventLoop) PostTask(task *Task) {
	if task.DelayedRunTime.After(time.Now()) {
		el.addToDelayedEventQueue(task)
	} else {
		el.addToPendingEventQueue(task)
	}
}

func (el *EventLoop) PostFunc(function func()) {
	var task = NewTask(function, 0)
	el.PostTask(task)
}

func (el *EventLoop) PostDelayedFunc(function func(), delayed time.Duration) {
	var task = NewTask(function, delayed)
	el.PostTask(task)
}

func (el *EventLoop) RepeatFunc(function func(stop *bool), delayed time.Duration) {
	var repeat *Task
	repeat = NewTask(func() {
		repeat.DelayedRunTime = time.Now().Add(delayed)
		var stop bool
		function(&stop)
		if stop {
			return
		}
		el.PostTask(repeat)
	}, 0)

	el.PostTask(repeat)
}

func (el *EventLoop) ShouldQuit() {
	el.shouldQuit = true
}

func (el *EventLoop) addToPendingEventQueue(task *Task) {
	var wasEmpty = (el.pendingTaskQueue.Len() == 0)
	el.pendingTaskQueue.PushBack(task)
	el.scheduleEvent(wasEmpty)
}

func (el *EventLoop) addToDelayedEventQueue(task *Task) {
	el.delayedTaskQueue.Push(task)
	var nextDelayedRunTime = el.delayedTaskQueue.Top().(*Task).DelayedRunTime
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

func (el *EventLoop) processTask() (handled bool) {
	var now = time.Now()

	for el.pendingTaskQueue.Len() != 0 {
		el.pendingTaskQueueMutex.Lock()
		var elem = el.pendingTaskQueue.Front()
		var task = elem.Value.(*Task)
		el.pendingTaskQueue.Remove(elem)
		el.pendingTaskQueueMutex.Unlock()

		if !task.DelayedRunTime.After(now) {
			task.Function()
			return el.pendingTaskQueue.Len() != 0
		} else {
			el.delayedTaskQueue.Push(task)
			if el.delayedTaskQueue.Top() == Interface(task) {
				el.pump.scheduleDelayedEvent(el.delayedTaskQueue.Top().(*Task).DelayedRunTime)
			}
		}
	}

	return false
}

func (el *EventLoop) processDelayedTask() (handled bool, nextDelayedEventTime time.Time) {
	if el.delayedTaskQueue.Empty() {
		el.recentTime = time.Now()
		return false, el.recentTime
	}

	// When we "fall behind", there will be a lot of tasks in the delayed event
	// queue that are ready to run. To increase efficiency when we fall behind,
	// we will only call Time.Now() intermittently, and then process all tasks
	// that are ready to run before calling it again. As a result, the more we
	// fall behind (and have a lot of ready-to-run delayed events), the more
	// efficiency we'll be at handling the events.
	var nextRunTime = el.delayedTaskQueue.Top().(*Task).DelayedRunTime
	if nextRunTime.After(el.recentTime) {
		el.recentTime = time.Now()
		if nextRunTime.After(el.recentTime) {
			return false, nextRunTime
		}
	}

	var event = el.delayedTaskQueue.Top().(*Task)
	el.delayedTaskQueue.Pop()

	event.Function()

	if !el.delayedTaskQueue.Empty() {
		nextDelayedEventTime = el.delayedTaskQueue.Top().(*Task).DelayedRunTime
	}

	return true, nextDelayedEventTime
}

func (el *EventLoop) Run() {
	el.pump.run()
}
