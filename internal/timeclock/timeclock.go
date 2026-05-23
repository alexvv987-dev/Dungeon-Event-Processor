package timeclock

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

var zeroPoint = time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)

func Parse(s string) (time.Time, error) {
	parts := strings.Split(s, ":")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", s)
	}
	h, e1 := strconv.Atoi(parts[0])
	m, e2 := strconv.Atoi(parts[1])
	sec, e3 := strconv.Atoi(parts[2])
	if e1 != nil || e2 != nil || e3 != nil {
		return time.Time{}, fmt.Errorf("invalid time format: %s", s)
	}
	if h < 0 || m < 0 || m > 59 || sec < 0 || sec > 59 {
		return time.Time{}, fmt.Errorf("invalid time format: %s", s)
	}
	return zeroPoint.Add(time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(sec)*time.Second), nil
}

func FormatTime(t time.Time) string {
	return fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
}

func FormatDuration(d time.Duration) string {
	if d < 0 {
		d = 0
	}
	h := int(d.Hours())
	m := int(d.Minutes()) % 60
	s := int(d.Seconds()) % 60
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func ZeroPoint() time.Time {
	return zeroPoint
}
