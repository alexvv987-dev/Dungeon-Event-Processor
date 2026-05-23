package event

import (
	"strings"
	"testing"

	"impulse/internal/timeclock"
)

func TestParseLine(t *testing.T) {
	ev, ok, err := ParseLine("[14:10:00] 2 11 60")
	if err != nil || !ok {
		t.Fatalf("parse: ok=%v err=%v", ok, err)
	}
	if ev.PlayerID != 2 || ev.ID != 11 || ev.Extra != "60" {
		t.Errorf("got %+v", ev)
	}
	if timeclock.FormatTime(ev.Time) != "14:10:00" {
		t.Errorf("time = %s", timeclock.FormatTime(ev.Time))
	}

	_, ok, _ = ParseLine("")
	if ok {
		t.Error("empty line should be skipped")
	}

	ev2, ok, _ := ParseLine("[14:00:00] 1 9 out of mana")
	if !ok || ev2.Extra != "out of mana" || ev2.ID != 9 {
		t.Errorf("multi-word extra: %+v", ev2)
	}
}

func TestReadAllMonotonic(t *testing.T) {
	_, err := ReadAll(strings.NewReader("[14:00:00] 1 1\n[13:00:00] 1 2\n"))
	if err == nil {
		t.Fatal("expected monotonic time error")
	}
}

func TestReadAllOK(t *testing.T) {
	got, err := ReadAll(strings.NewReader("[14:00:00] 1 1\n[14:01:00] 1 2\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d, want 2", len(got))
	}
}
