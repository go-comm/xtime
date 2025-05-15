package xtime

import (
	"context"
	"fmt"
	"sync/atomic"
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

func SetInterval(ctx context.Context, period time.Duration, f func() error) Future {
	f = ErrRecovered(f)
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	go func() {
		t := AcquireTicker(nil, period)
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

func setTimeout(ctx context.Context, delay time.Duration, f func() error) Future {
	run := funcRunner(ctx, f)
	if delay <= time.Millisecond {
		go run()
		return FutureFunc(func() {})
	}
	t := time.AfterFunc(delay, run)
	return FutureFunc(func() { ReleaseTimer(t) })
}

func SetTimeout(ctx context.Context, delay time.Duration, f func() error) Future {
	return setTimeout(ctx, delay, ErrRecovered(f))
}

func SetSchedule(ctx context.Context, delay time.Duration, period time.Duration, f func() error) Future {
	f = ErrRecovered(f)
	var h func() error
	var future Future
	var done int32

	h = func() error {
		if atomic.LoadInt32(&done) == 1 {
			return nil
		}
		t0 := time.Now()
		err := f()
		if atomic.LoadInt32(&done) == 1 {
			return err
		}
		t1 := time.Now()
		delta := t1.Sub(t0)
		if period > 0 {
			d := (delta+period)/period*period - delta
			future = setTimeout(ctx, d, h)
		}
		return err
	}

	future = setTimeout(ctx, delay, h)
	return FutureFunc(func() {
		atomic.StoreInt32(&done, 1)
		future.Stop()
	})
}

func Sleep(ctx context.Context, d time.Duration) error {
	var err error
	if d <= 0 {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break
		default:
		}
		return err
	}

	var t *time.Timer
	t = AcquireTimer(t, d)
	defer ReleaseTimer(t)
LOOP:
	for {
		select {
		case <-ctx.Done():
			err = ctx.Err()
			break LOOP
		case <-t.C:
			break LOOP
		}
	}
	return err
}

func funcRunner(ctx context.Context, f func() error) func() {
	return func() {
		select {
		case <-ctx.Done():
			return
		default:
			f()
		}
	}
}
