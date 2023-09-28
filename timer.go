package xtime

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

func id() int64 {
	n := rand.Int63()
	if n == 0 {
		for n != 0 {
			n = rand.Int63()
		}
	}
	if n < 0 {
		n = -n
	}
	return n >> 10 //  for javascript
}

func now() int64 {
	return time.Now().UnixNano()
}

func when(d time.Duration) int64 {
	if d <= 0 {
		return now()
	}
	return now() + int64(d)
}

type TimerOption func(e *timerEntry)

func Handler(h func(interface{})) TimerOption { return func(e *timerEntry) { e.h = h } }

type TimerFuture interface {
	ID() int64
	Stop() bool
}

type timerFuture struct {
	id int64
	m  *TimerManager
}

func (t *timerFuture) ID() int64 {
	return t.id
}

func (t *timerFuture) Stop() bool {
	return t.m.delByID(t.id)
}

type timerEntry struct {
	id     int64
	when   int64
	period int64
	data   interface{}
	h      func(interface{})

	m     *TimerManager
	timer *time.Timer
}

func (e *timerEntry) Do() {
	m := e.m
	h := e.h
	data := e.data

	if m == nil {
		return
	}

	if e.period > 0 {
		o := &timerEntry{id: e.id, when: e.when + e.period, period: e.period, data: data, h: h}
		m.add(o, true)
	} else {
		m.del(e)
	}
	if h == nil {
		h = m.h
	}
	if h != nil {
		h(data)
	}
}

func (e *timerEntry) Reset() {
	e.timer = ReleaseTimer(e.timer)
	d := time.Duration(e.when - now())
	if d <= 0 {
		go e.Do()
		return
	}
	e.timer = time.AfterFunc(d, e.Do)
}

func (e *timerEntry) Future() TimerFuture {
	return &timerFuture{id: e.id, m: e.m}
}

func releaseTimerEntry(e *timerEntry) {
	if e != nil {
		e.timer = ReleaseTimer(e.timer)
		e.data = nil
		e.h = nil
		e.m = nil
	}
}

type refTimerEntry struct {
	p unsafe.Pointer // *timerEntry
}

func (ref *refTimerEntry) Ref() *timerEntry {
	p := atomic.LoadPointer(&ref.p)
	return (*timerEntry)(p)
}

func (ref *refTimerEntry) Set(e *timerEntry) {
	atomic.StorePointer(&ref.p, unsafe.Pointer(e))
}

func NewTimerManager() *TimerManager {
	m := &TimerManager{}
	return m
}

type TimerManager struct {
	h         func(interface{})
	timerlist sync.Map
}

func (m *TimerManager) SetHandler(h func(interface{})) {
	m.h = h
}

func (m *TimerManager) add(e *timerEntry, cover bool, opts ...TimerOption) bool {
	e.m = m

	if len(opts) > 0 {
		for _, opt := range opts {
			opt(e)
		}
	}

	o, found := m.timerlist.LoadOrStore(e.id, &refTimerEntry{p: unsafe.Pointer(e)})
	if found && !cover {
		return false
	}
	ref := o.(*refTimerEntry)
	if found && cover {
		olde := ref.Ref()
		ref.Set(e)
		releaseTimerEntry(olde)
	}
	e.Reset()
	return true
}

func (m *TimerManager) del(e *timerEntry) bool {
	defer releaseTimerEntry(e)
	return m.delByID(e.id)
}

func (m *TimerManager) delByID(id int64) bool {
	o, ok := m.timerlist.LoadAndDelete(id)
	if ok {
		ref := o.(*refTimerEntry)
		if ref != nil {
			releaseTimerEntry(ref.Ref())
		}
		return true
	}
	return false
}

func (m *TimerManager) Has(id int64) bool {
	_, f := m.timerlist.Load(id)
	return f
}

func (m *TimerManager) Put(id int64, delay time.Duration, period time.Duration, x interface{}, opts ...TimerOption) TimerFuture {
	e := &timerEntry{id: id, when: when(delay), period: int64(period), data: x}
	m.add(e, true, opts...)
	return e.Future()
}

func (m *TimerManager) PutAt(id int64, t time.Time, period time.Duration, x interface{}, opts ...TimerOption) TimerFuture {
	e := &timerEntry{id: id, when: t.UnixNano(), period: int64(period), data: x}
	m.add(e, true, opts...)
	return e.Future()
}

func (m *TimerManager) post(e *timerEntry, opts ...TimerOption) TimerFuture {
	e.id = id()
	for !m.add(e, false, opts...) {
		e.id = id()
	}
	return e.Future()
}

func (m *TimerManager) PostDelay(delay time.Duration, x interface{}, opts ...TimerOption) TimerFuture {
	e := &timerEntry{when: when(delay), data: x}
	return m.post(e, opts...)
}

func (m *TimerManager) PostAt(t time.Time, x interface{}, opts ...TimerOption) TimerFuture {
	e := &timerEntry{when: t.UnixNano(), data: x}
	return m.post(e, opts...)
}

func (m *TimerManager) PostPeriod(delay time.Duration, period time.Duration, x interface{}, opts ...TimerOption) TimerFuture {
	e := &timerEntry{when: when(delay), period: int64(period), data: x}
	return m.post(e, opts...)
}

func (m *TimerManager) StopByID(id int64) {
	m.delByID(id)
}

func (m *TimerManager) Stop(future TimerFuture) {
	future.Stop()
}
