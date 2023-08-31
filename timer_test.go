package xtime

import (
	"testing"
	"time"
)

func TestPostDelay(t *testing.T) {
	var m TimerManager
	stop := make(chan bool)

	m.SetHandler(func(x interface{}) {
		i := x.(int)
		if i >= 0 {
			future := m.PostDelay(1*time.Second, i-1)
			_ = future

		} else {
			stop <- true
		}
	})

	m.PostDelay(0, 10)
	<-stop
}

func TestPostPeriod(t *testing.T) {
	var m TimerManager
	stop := make(chan bool)

	var c int
	var future TimerFuture

	m.SetHandler(func(x interface{}) {
		i := x.(int)
		_ = i
		c++

		if c > 10 {
			future.Stop()
			stop <- true
		}

	})

	future = m.PostPeriod(time.Second, time.Second, 10)
	<-stop
}
