package task

import (
	"library"
	"time"
)

type Test struct {
}

func (t *Test) Start(complete chan<- int) {
	//complete <- 0
	i := 0
	for {
		i++
		library.Println("log/setLeaveInStudent", i)
		if i > 100 {
			complete <- 0
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
}
