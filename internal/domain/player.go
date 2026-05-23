package domain

import (
	"time"

	"impulse/internal/config"
)

type TrialState string

const (
	StateSuccess TrialState = "SUCCESS"
	StateFail    TrialState = "FAIL"
	StateDisqual TrialState = "DISQUAL"
)

const MaxHealth = 100

type Player struct {
	ID         int
	Registered bool
	InDungeon  bool
	Dead       bool
	Disqual    bool
	Health     int

	EnteredAt   time.Time
	LeftAt      time.Time
	DungeonDone bool

	CurrentFloor   int
	OnBossFloor    bool
	MonstersLeft   []int
	FloorsCleared  []bool
	FloorEnteredAt time.Time
	FloorTimeSpent []time.Duration

	BossFloorEnteredAt time.Time
	BossKillTime       time.Duration
	BossDefeated       bool
}

func NewPlayer(id int, regularFloors int, cfg config.Config) *Player {
	mobs := make([]int, regularFloors)
	cleared := make([]bool, regularFloors)
	spent := make([]time.Duration, regularFloors)
	for i := range mobs {
		mobs[i] = cfg.Monsters
	}
	return &Player{
		ID:             id,
		Health:         MaxHealth,
		MonstersLeft:   mobs,
		FloorsCleared:  cleared,
		FloorTimeSpent: spent,
	}
}

func (p *Player) IsActive() bool {
	return p.InDungeon && !p.DungeonDone && !p.Dead && !p.Disqual
}

func (p *Player) AllRegularFloorsCleared() bool {
	for _, done := range p.FloorsCleared {
		if !done {
			return false
		}
	}
	return len(p.FloorsCleared) > 0
}

func (p *Player) TryFloorClear(regularFloors int, t time.Time) {
	if p.OnBossFloor || p.CurrentFloor == 0 || p.CurrentFloor > regularFloors {
		return
	}
	idx := p.CurrentFloor - 1
	if !p.FloorsCleared[idx] && p.MonstersLeft[idx] == 0 {
		p.FloorsCleared[idx] = true
		p.FloorTimeSpent[idx] += t.Sub(p.FloorEnteredAt)
		p.FloorEnteredAt = t
	}
}

func (p *Player) StopFloorTimer(regularFloors int, t time.Time) {
	if p.OnBossFloor || p.CurrentFloor == 0 || p.CurrentFloor > regularFloors {
		return
	}
	idx := p.CurrentFloor - 1
	if !p.FloorsCleared[idx] {
		p.FloorTimeSpent[idx] += t.Sub(p.FloorEnteredAt)
	}
}

func (p *Player) RestoreHealth(amount int) int {
	p.Health += amount
	if p.Health > MaxHealth {
		p.Health = MaxHealth
	}
	return amount
}

func (p *Player) ApplyDamage(dmg int) {
	p.Health -= dmg
	if p.Health <= 0 {
		p.Health = 0
	}
}

func (p *Player) MarkDead(t time.Time) {
	p.Dead = true
	p.DungeonDone = true
	p.LeftAt = t
}

func (p *Player) EndChallenge(t time.Time) {
	p.DungeonDone = true
	p.LeftAt = t
}

func (p *Player) FinalState() TrialState {
	if !p.Registered || p.Disqual {
		return StateDisqual
	}
	if p.BossDefeated && p.AllRegularFloorsCleared() {
		return StateSuccess
	}
	return StateFail
}

func (p *Player) DungeonDuration() time.Duration {
	if p.EnteredAt.IsZero() || p.LeftAt.IsZero() {
		return 0
	}
	return p.LeftAt.Sub(p.EnteredAt)
}

func (p *Player) AvgFloorTime(regularFloors int) time.Duration {
	if regularFloors <= 0 {
		return 0
	}
	var total time.Duration
	count := 0
	for _, d := range p.FloorTimeSpent {
		if d > 0 {
			total += d
			count++
		}
	}
	if count == 0 {
		return 0
	}
	return total / time.Duration(count)
}
