package xtime

import (
	"context"
	"time"
)

func SetExecute(ctx context.Context, retries int, period time.Duration, f func() (continued bool, err error)) error {
	var err error
	var continued bool
	var cnt int
	for ; retries < 0 || cnt <= retries; cnt++ {
		if cnt > 0 {
			if err = Sleep(ctx, period); err != nil {
				break
			}
		}
		continued, err = f()
		if continued {
			continue
		}
		if err != nil {
			break
		}
	}
	return err
}

func SetRetry(ctx context.Context, retries int, f func() (continued bool, err error)) error {
	return SetExecute(ctx, retries, 0, f)
}

func SetUtil(ctx context.Context, period time.Duration, f func() (continued bool, err error)) error {
	return SetExecute(ctx, -1, period, f)
}

func SetBackoff(ctx context.Context, min, max time.Duration, f func() (continued bool, err error)) error {
	var err error
	var continued bool
	if min < 0 {
		min = 1 * time.Second
	}
	if max < 0 {
		max = 3 * time.Second
	}
	if min > max {
		min, max = max, min
	}
	backoff := min
	var cnt int
	for ; ; cnt++ {
		if cnt > 0 {
			if err = Sleep(ctx, backoff); err != nil {
				break
			}
			backoff = backoff * 9 / 8
			if backoff > max {
				backoff = min
			}
		}
		continued, err = f()
		if continued {
			continue
		}
		if err != nil {
			break
		}
	}
	return err
}
