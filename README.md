# evloop
a eventloop for golang. We can write schedule use evloop. or as async eventloop thread during developing.

# sample code.

```
package main

import (
	"fmt"
	"time"

	"github.com/sj20082663/evloop"
)

var loop *evloop.EventLoop

func main() {
	loop = evloop.NewEventLoop()
	thread()
	loop.Run()
}

func thread() {
	loop.PostFunc(func() {
		fmt.Printf("1. a coin has two side\n")
	})

	loop.PostDelayedFunc(func() {
		fmt.Printf("4. a bird in the hand is worth two in the bush\n")
	}, 3*time.Second)

	loop.PostDelayedFunc(func() {
		fmt.Printf("3. There's always time, Time is first\n")
	}, 2*time.Second)

	loop.RepeatFunc(func(stop *bool) {
		fmt.Printf("2. too be or not to be is a question\n")
	}, time.Second)
}
```

# Copyright
Copyright 2015 By Jash. All rights reserved.
Use of this source code is governed by a BSD-style
license that can be found in the LICENSE file.

# Contact
mail: shijian0912@163.com
