package xtime

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func ErrRecovered(f func() error) func() error {
	return func() (err error) {
		defer func() {
			if ierr := recover(); ierr != nil {
				var ok bool
				if err, ok = ierr.(error); !ok {
					err = fmt.Errorf("%v", ierr)
				}
			}
		}()
		return f()
	}
}

func AcquireTimer(t *time.Timer, d time.Duration) *time.Timer {
	if t == nil {
		t = time.NewTimer(d)
	} else {
		ReleaseTimer(t)
		t.Reset(d)
	}
	return t
}

func ReleaseTimer(t *time.Timer) *time.Timer {
	if t != nil && !t.Stop() {
		select {
		case <-t.C:
		default:
		}
	}
	return t
}

func AcquireTicker(t *time.Ticker, d time.Duration) *time.Ticker {
	if t == nil {
		t = time.NewTicker(d)
	} else {
		ReleaseTicker(t)
		t.Reset(d)
	}
	return t
}

func ReleaseTicker(t *time.Ticker) *time.Ticker {
	if t != nil {
		t.Stop()
		select {
		case <-t.C:
		default:
		}
	}
	return t
}

type Future interface {
	Stop()
}

type tickerFuture struct {
	mutex sync.Mutex
	t     *time.Ticker
}

func (f *tickerFuture) Stop() {
	f.mutex.Lock()
	t := f.t
	f.t = nil
	f.mutex.Unlock()
	ReleaseTicker(t)
}

func SetInterval(ctx context.Context, d time.Duration, f func() error) Future {
	f = ErrRecovered(f)
	future := &tickerFuture{t: AcquireTicker(nil, d)}
	go func() {
		defer future.Stop()
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-future.t.C:
				f()
			}
		}
	}()
	return future
}

type timeoutFuture struct {
	mutex sync.Mutex
	t     *time.Timer
}

func (f *timeoutFuture) Stop() {
	f.mutex.Lock()
	t := f.t
	f.t = nil
	f.mutex.Unlock()
	ReleaseTimer(t)
}

func SetTimeout(ctx context.Context, d time.Duration, f func() error) Future {
	f = ErrRecovered(f)
	var future = &timeoutFuture{t: AcquireTimer(nil, d)}
	go func() {
		defer future.Stop()
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-future.t.C:
				f()
				break LOOP
			}
		}
	}()
	return future
}

func WithTimeout(ctx context.Context, d time.Duration, f func(ctx context.Context) error) func() error {
	return func() error {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, d)
		defer cancel()
		return f(ctx)
	}
}

func Sleep(ctx context.Context, d time.Duration) {
	var t *time.Timer
	t = AcquireTimer(t, d)
	defer ReleaseTimer(t)
LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-t.C:
			break LOOP
		}
	}
}
