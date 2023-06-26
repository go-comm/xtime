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

func SetInterval(ctx context.Context, d time.Duration, f func() error) {
	f = ErrRecovered(f)
	go func() {
		t := time.NewTicker(d)
		defer t.Stop()
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
}

func SetTimeout(ctx context.Context, d time.Duration, f func() error) {
	f = ErrRecovered(f)
	go func() {
		var t *time.Timer
		t = AcquireTimer(t, d)
		defer ReleaseTimer(t)

	LOOP:
		for {
			select {
			case <-ctx.Done():
				break LOOP
			case <-t.C:
				f()
				break LOOP
			}
		}

	}()
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
