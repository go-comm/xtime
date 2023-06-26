package xtime

import (
	"time"
)

var (
	Layout      = "01/02 03:04:05PM '06 -0700"
	ANSIC       = "Mon Jan _2 15:04:05 2006"
	UnixDate    = "Mon Jan _2 15:04:05 MST 2006"
	RubyDate    = "Mon Jan 02 15:04:05 -0700 2006"
	RFC822      = "02 Jan 06 15:04 MST"
	RFC822Z     = "02 Jan 06 15:04 -0700"
	RFC850      = "Monday, 02-Jan-06 15:04:05 MST"
	RFC1123     = "Mon, 02 Jan 2006 15:04:05 MST"
	RFC1123Z    = "Mon, 02 Jan 2006 15:04:05 -0700"
	RFC3339     = "2006-01-02T15:04:05Z07:00"
	RFC3339Nano = "2006-01-02T15:04:05.999999999Z07:00"
	ISO8601     = "2006-01-02T15:04:05Z"
	DateTime    = "2006-01-02 15:04:05"
	DateOnly    = "2006-01-02"
	TimeOnly    = "15:04:05"
)

var (
	LocationAsiaShanghai = "Asia/Shanghai"
)

func MustLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return loc
}

var gLocal *time.Location = time.Local

func Local() *time.Location {
	return gLocal
}

func SetLocal(loc *time.Location) *time.Location {
	gLocal = loc
	return loc
}

func Now() time.Time {
	return time.Now().In(Local())
}

func Compare(t time.Time, u time.Time) int {
	tc := t.UnixNano()
	uc := u.UnixNano()
	if tc > uc {
		return 1
	}
	if tc < uc {
		return -1
	}
	return 0
}

func Less(t time.Time, u time.Time) bool {
	return Compare(t, u) < 0
}

func Equal(t time.Time, u time.Time) bool {
	return Compare(t, u) == 0
}

func Sub(t time.Time, u time.Time) time.Duration {
	return t.Sub(u)
}

func StartOfYear(t time.Time) time.Time {
	year, _, _ := t.Date()
	return time.Date(year, 1, 1, 0, 0, 0, 0, t.Location())
}

func StartOfMonth(t time.Time) time.Time {
	year, month, _ := t.Date()
	return time.Date(year, month, 1, 0, 0, 0, 0, t.Location())
}

func StartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func StartOfHour(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, _, _ := t.Clock()
	return time.Date(year, month, day, hour, 0, 0, 0, t.Location())
}

func StartOfMinute(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, _ := t.Clock()
	return time.Date(year, month, day, hour, min, 0, 0, t.Location())
}

func StartOfSecond(t time.Time) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(year, month, day, hour, min, sec, 0, t.Location())
}

func EndOfYear(t time.Time) time.Time {
	return AddYear(StartOfYear(t), 1).Add(-time.Nanosecond)
}

func EndOfMonth(t time.Time) time.Time {
	return AddMonth(StartOfMonth(t), 1).Add(-time.Nanosecond)
}

func EndOfDay(t time.Time) time.Time {
	return StartOfDay(t).Add(time.Hour * 24).Add(-time.Nanosecond)
}

func EndOfHour(t time.Time) time.Time {
	return StartOfHour(t).Add(time.Hour).Add(-time.Nanosecond)
}

func EndOfMinute(t time.Time) time.Time {
	return StartOfMinute(t).Add(time.Minute).Add(-time.Nanosecond)
}

func EndOfSecond(t time.Time) time.Time {
	return StartOfSecond(t).Add(time.Second).Add(-time.Nanosecond)
}

func SinceYear(t time.Time) time.Duration {
	return t.Sub(StartOfYear(t))
}

func SinceMonth(t time.Time) time.Duration {
	return t.Sub(StartOfMonth(t))
}

func SinceDay(t time.Time) time.Duration {
	return t.Sub(StartOfDay(t))
}

func SinceHour(t time.Time) time.Duration {
	return t.Sub(StartOfHour(t))
}

func SinceMinute(t time.Time) time.Duration {
	return t.Sub(StartOfMinute(t))
}

func SinceSecord(t time.Time) time.Duration {
	return t.Sub(StartOfSecond(t))
}

func AddYear(t time.Time, y int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(year+y, month, day, hour, min, sec, t.Nanosecond(), t.Location())
}

func AddMonth(t time.Time, m int) time.Time {
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	return time.Date(year, month+time.Month(m), day, hour, min, sec, t.Nanosecond(), t.Location())
}

func AddDay(t time.Time, d int) time.Time {
	return t.Add(time.Duration(d) * time.Hour * 24)
}

func AddHour(t time.Time, h int) time.Time {
	return t.Add(time.Duration(h) * time.Hour)
}

func AddMinute(t time.Time, m int) time.Time {
	return t.Add(time.Duration(m) * time.Minute)
}

func AddSecond(t time.Time, s int) time.Time {
	return t.Add(time.Duration(s) * time.Second)
}

func SetClock(t time.Time, hour, min, sec int) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, hour, min, sec, 0, t.Location())
}

func SetDate(t time.Time, year int, month time.Month, day int) time.Time {
	hour, min, sec := t.Clock()
	return time.Date(year, month, day, hour, min, sec, t.Nanosecond(), t.Location())
}

func Every(from time.Time, to time.Time, interval time.Duration, f func(u time.Time, off int)) {
	var i int
	for u := from; Compare(u, to) <= 0; u = u.Add(interval) {
		f(u, i)
		i++
	}
}

func EveryDay(from time.Time, to time.Time, f func(u time.Time, off int)) {
	Every(from, to, time.Hour*24, f)
}

func EveryHour(from time.Time, to time.Time, f func(u time.Time, off int)) {
	Every(from, to, time.Hour, f)
}

func EveryMinute(from time.Time, to time.Time, f func(u time.Time, off int)) {
	Every(from, to, time.Minute, f)
}

func EverySecond(from time.Time, to time.Time, f func(u time.Time, off int)) {
	Every(from, to, time.Second, f)
}
