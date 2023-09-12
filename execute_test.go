package xtime

import (
	"context"
	"fmt"
	"testing"
)

func TestSetRetry(t *testing.T) {

	var maxCnt int = 5
	var cur int = 0
	err := SetRetry(context.TODO(), maxCnt-1, func() (continued bool, err error) {
		cur++
		if cur < maxCnt {
			return true, fmt.Errorf("no match, cur:%v, max:%v", cur, maxCnt)
		}
		return false, nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if cur != maxCnt {
		t.Fatal(cur, maxCnt)
	}

}
