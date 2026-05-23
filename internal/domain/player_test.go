package domain

import (
	"testing"
	"time"

	"impulse/internal/config"
	"impulse/internal/timeclock"
)

func TestNewPlayerInitialization(t *testing.T) {
	cfg := config.Config{Floors: 4, Monsters: 3, OpenAt: "10:00:00", Duration: 2}
	floors := cfg.RegularFloors()
	p := NewPlayer(42, floors, cfg)

	if p.ID != 42 {
		t.Errorf("ID = %d, want 42", p.ID)
	}
	if p.Health != MaxHealth {
		t.Errorf("Health = %d, want %d", p.Health, MaxHealth)
	}
	if len(p.MonstersLeft) != floors {
		t.Errorf("MonstersLeft len = %d, want %d", len(p.MonstersLeft), floors)
	}
	for i, m := range p.MonstersLeft {
		if m != cfg.Monsters {
			t.Errorf("MonstersLeft[%d] = %d, want %d", i, m, cfg.Monsters)
		}
	}
	if p.Registered || p.InDungeon || p.Dead || p.Disqual {
		t.Error("new player should have all bool flags false")
	}
}

func TestRestoreHealthCap(t *testing.T) {
	p := &Player{Health: 90}
	p.RestoreHealth(80)
	if p.Health != 100 {
		t.Errorf("Health = %d, want 100", p.Health)
	}
	p.Health = 20
	p.RestoreHealth(15)
	if p.Health != 35 {
		t.Errorf("Health = %d, want 35", p.Health)
	}
}

func TestAvgFloorTimeOnlyCleared(t *testing.T) {
	p := &Player{
		FloorTimeSpent: []time.Duration{5 * time.Minute, 0, 10 * time.Minute},
		FloorsCleared:  []bool{true, false, true},
	}
	got := p.AvgFloorTime(3)
	want := 7*time.Minute + 30*time.Second
	if got != want {
		t.Errorf("AvgFloorTime = %v, want %v", got, want)
	}
}

func TestFinalState(t *testing.T) {
	p := &Player{Registered: true, BossDefeated: true, FloorsCleared: []bool{true}}
	if p.FinalState() != StateSuccess {
		t.Errorf("got %s, want SUCCESS", p.FinalState())
	}
	p.BossDefeated = false
	if p.FinalState() != StateFail {
		t.Errorf("got %s, want FAIL", p.FinalState())
	}
	p.Disqual = true
	if p.FinalState() != StateDisqual {
		t.Errorf("got %s, want DISQUAL", p.FinalState())
	}
}

func TestDungeonDuration(t *testing.T) {
	entered, _ := timeclock.Parse("14:40:00")
	left, _ := timeclock.Parse("15:04:00")
	p := &Player{EnteredAt: entered, LeftAt: left}
	if got := timeclock.FormatDuration(p.DungeonDuration()); got != "00:24:00" {
		t.Errorf("duration = %s, want 00:24:00", got)
	}
}
