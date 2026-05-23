package timeclock

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	cases := []struct {
		input   string
		wantErr bool
		h, m, s int
	}{
		{"14:00:00", false, 14, 0, 0},
		{"09:05:02", false, 9, 5, 2},
		{"00:00:00", false, 0, 0, 0},
		{"bad", true, 0, 0, 0},
		{"14:00", true, 0, 0, 0},
		{"14:60:00", true, 0, 0, 0},
	}
	for _, tc := range cases {
		got, err := Parse(tc.input)
		if tc.wantErr {
			if err == nil {
				t.Errorf("Parse(%q) expected error", tc.input)
			}
			continue
		}
		if err != nil {
			t.Errorf("Parse(%q) unexpected error: %v", tc.input, err)
			continue
		}
		want := ZeroPoint().Add(time.Duration(tc.h)*time.Hour + time.Duration(tc.m)*time.Minute + time.Duration(tc.s)*time.Second)
		if !got.Equal(want) {
			t.Errorf("Parse(%q) = %v, want %v", tc.input, FormatTime(got), FormatTime(want))
		}
	}
}

func TestFormatDuration(t *testing.T) {
	cases := []struct {
		d    time.Duration
		want string
	}{
		{0, "00:00:00"},
		{time.Hour + 30*time.Minute + 5*time.Second, "01:30:05"},
		{24 * time.Hour, "24:00:00"},
		{-time.Minute, "00:00:00"},
		{11 * time.Minute, "00:11:00"},
	}
	for _, tc := range cases {
		if got := FormatDuration(tc.d); got != tc.want {
			t.Errorf("FormatDuration(%v) = %q, want %q", tc.d, got, tc.want)
		}
	}
}

func TestFormatTimeRoundTrip(t *testing.T) {
	for _, input := range []string{"08:03:07", "14:00:00", "00:00:01"} {
		parsed, err := Parse(input)
		if err != nil {
			t.Fatalf("Parse(%q): %v", input, err)
		}
		if got := FormatTime(parsed); got != input {
			t.Errorf("FormatTime(Parse(%q)) = %q", input, got)
		}
	}
}
