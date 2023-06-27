package xtime

import (
	"context"
	"fmt"
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

type FutureFunc func()

func (f FutureFunc) Stop() { f() }

func ClearFuture(f Future) {
	if f != nil {
		f.Stop()
	}
}

func SetInterval(ctx context.Context, d time.Duration, f func() error) Future {
	f = ErrRecovered(f)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		t := AcquireTicker(nil, d)
		defer ReleaseTicker(t)
	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-t.C:
				f()
			}
		}
	}()
	return FutureFunc(func() { cancel() })
}

func SetTimeout(ctx context.Context, d time.Duration, f func() error) Future {
	f = ErrRecovered(f)
	t := time.AfterFunc(d, func() {
		select {
		case <-ctx.Done():
			return
		default:
			f()
		}
	})
	return FutureFunc(func() { ReleaseTimer(t) })
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
