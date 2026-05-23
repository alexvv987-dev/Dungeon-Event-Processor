package domain

import (
	"fmt"
	"time"

	"impulse/internal/config"
	"impulse/internal/timeclock"
)

type Schedule struct {
	OpenAt  time.Time
	CloseAt time.Time
}

func NewSchedule(cfg config.Config) (Schedule, error) {
	openAt, err := timeclock.Parse(cfg.OpenAt)
	if err != nil {
		return Schedule{}, fmt.Errorf("bad OpenAt: %w", err)
	}
	closeAt := openAt.Add(time.Duration(cfg.Duration) * time.Hour)
	return Schedule{OpenAt: openAt, CloseAt: closeAt}, nil
}

func (s Schedule) IsOpen(t time.Time) bool {
	return !t.Before(s.OpenAt) && t.Before(s.CloseAt)
}

func (s Schedule) CloseActivePlayers(players map[int]*Player, t time.Time) {
	if t.Before(s.CloseAt) {
		return
	}
	for _, p := range players {
		if p.InDungeon && !p.DungeonDone {
			p.DungeonDone = true
			p.LeftAt = s.CloseAt
		}
	}
}

func CloseRemainingPlayers(players map[int]*Player, closeAt time.Time) {
	for _, p := range players {
		if p.InDungeon && !p.DungeonDone {
			p.DungeonDone = true
			p.LeftAt = closeAt
		}
	}
}
