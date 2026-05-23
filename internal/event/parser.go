package event

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"impulse/internal/timeclock"
)

const (
	Register       = 1
	EnterDungeon   = 2
	KillMonster    = 3
	NextFloor      = 4
	PrevFloor      = 5
	EnterBoss      = 6
	KillBoss       = 7
	LeaveDungeon   = 8
	CannotContinue = 9
	RestoreHealth  = 10
	ReceiveDamage  = 11
)

const (
	OutDisqualified   = 31
	OutDead           = 32
	OutImpossibleMove = 33
)

type Event struct {
	Time     time.Time
	PlayerID int
	ID       int
	Extra    string
}

func ParseLine(line string) (ev Event, ok bool, err error) {
	line = strings.TrimSpace(line)
	if line == "" || !strings.HasPrefix(line, "[") {
		return Event{}, false, nil
	}
	bracket := strings.Index(line, "]")
	if bracket < 0 {
		return Event{}, false, nil
	}
	t, err := timeclock.Parse(line[1:bracket])
	if err != nil {
		return Event{}, false, nil
	}
	fields := strings.Fields(strings.TrimSpace(line[bracket+1:]))
	if len(fields) < 2 {
		return Event{}, false, nil
	}
	pid, e1 := strconv.Atoi(fields[0])
	eid, e2 := strconv.Atoi(fields[1])
	if e1 != nil || e2 != nil {
		return Event{}, false, nil
	}
	extra := ""
	if len(fields) > 2 {
		extra = strings.Join(fields[2:], " ")
	}
	return Event{Time: t, PlayerID: pid, ID: eid, Extra: extra}, true, nil
}

func ReadAll(r io.Reader) ([]Event, error) {
	var events []Event
	var prev time.Time
	seenFirst := false

	sc := bufio.NewScanner(r)
	for sc.Scan() {
		ev, ok, err := ParseLine(sc.Text())
		if err != nil {
			return nil, err
		}
		if !ok {
			continue
		}
		if seenFirst && ev.Time.Before(prev) {
			return nil, fmt.Errorf("event time %s is before previous event %s",
				timeclock.FormatTime(ev.Time), timeclock.FormatTime(prev))
		}
		prev = ev.Time
		seenFirst = true
		events = append(events, ev)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read events: %w", err)
	}
	return events, nil
}
