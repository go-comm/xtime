package xtime

import (
	"strings"
	"time"
)

func MustParse(layout string, value string) time.Time {
	t, err := ParseInLocation(layout, value, Local())
	if err != nil {
		panic(err)
	}
	return t
}

func Parse(layout string, value string) (t time.Time, err error) {
	return ParseInLocation(layout, value, Local())
}

func MustParseInLocation(layout string, value string, loc *time.Location) time.Time {
	t, err := ParseInLocation(layout, value, loc)
	if err != nil {
		panic(err)
	}
	return t
}

func ParseInLocation(layout string, value string, loc *time.Location) (t time.Time, err error) {
	if strings.Contains(value, "/") {
		value = strings.ReplaceAll(value, "/", "-")
	}
	if len(layout) > 0 {
		return time.ParseInLocation(layout, value, loc)
	}

	n := len(value)
	if n == len(DateTime) {
		if t, err = time.ParseInLocation(DateTime, value, loc); err == nil {
			return
		}
	}
	if n == len(DateOnly) {
		if t, err = time.ParseInLocation(DateOnly, value, loc); err == nil {
			return
		}
	}
	if n == len(RFC3339) {
		if t, err = time.ParseInLocation(RFC3339, value, loc); err == nil {
			return
		}
	}
	if n == len(RFC3339Nano) {
		if t, err = time.ParseInLocation(RFC3339Nano, value, loc); err == nil {
			return
		}
	}
	if n == len(ISO8601) {
		if t, err = time.ParseInLocation(ISO8601, value, loc); err == nil {
			return
		}
	}
	if n == len(RFC1123) {
		if t, err = time.ParseInLocation(RFC1123, value, loc); err == nil {
			return
		}
	}
	if n == len(RFC1123Z) {
		if t, err = time.ParseInLocation(RFC1123Z, value, loc); err == nil {
			return
		}
	}
	if n == len(RFC822) {
		if t, err = time.ParseInLocation(RFC822, value, loc); err == nil {
			return
		}
	}
	if n == len(RFC822Z) {
		if t, err = time.ParseInLocation(RFC822Z, value, loc); err == nil {
			return
		}
	}
	if n == len(RubyDate) {
		if t, err = time.ParseInLocation(RubyDate, value, loc); err == nil {
			return
		}
	}
	if n == len(UnixDate) {
		if t, err = time.ParseInLocation(UnixDate, value, loc); err == nil {
			return
		}
	}
	if n == len(ANSIC) {
		if t, err = time.ParseInLocation(ANSIC, value, loc); err == nil {
			return
		}
	}
	if n == len(Layout) {
		if t, err = time.ParseInLocation(Layout, value, loc); err == nil {
			return
		}
	}
	return time.ParseInLocation(layout, value, loc)
}
