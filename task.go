// Copyright 2015 By Jash. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
// mail: shijian0912@163.com

package evloop

import "time"

type Task struct {
	Function       func()
	DelayedRunTime time.Time
}

func NewTask(function func(), delayed time.Duration) *Task {
	var task = new(Task)
	task.Function = function
	task.DelayedRunTime = time.Now().Add(delayed)
	return task
}

func (e *Task) Precede(a Interface) bool {
	var ea = a.(*Task)
	return e.DelayedRunTime.Before(ea.DelayedRunTime)
}
