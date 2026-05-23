package engine

import (
	"log/slog"
	"strconv"

	"impulse/internal/domain"
	"impulse/internal/event"
	"impulse/internal/output"
)

func (s *session) handle(ev event.Event) {
	s.schedule.CloseActivePlayers(s.players, ev.Time)

	p := s.fetchPlayer(ev.PlayerID)
	active := p.IsActive()

	switch ev.ID {
	case event.Register:
		s.handleRegister(ev, p)
	case event.EnterDungeon:
		s.handleEnterDungeon(ev, p)
	case event.KillMonster:
		if !active {
			return
		}
		s.handleKillMonster(ev, p)
	case event.NextFloor:
		if !active {
			return
		}
		s.handleNextFloor(ev, p)
	case event.PrevFloor:
		if !active {
			return
		}
		s.handlePrevFloor(ev, p)
	case event.EnterBoss:
		if !active {
			return
		}
		s.handleEnterBoss(ev, p)
	case event.KillBoss:
		if !active {
			return
		}
		s.handleKillBoss(ev, p)
	case event.LeaveDungeon:
		if !active {
			return
		}
		s.handleLeaveDungeon(ev, p)
	case event.CannotContinue:
		if !active {
			return
		}
		s.handleCannotContinue(ev, p)
	case event.RestoreHealth:
		if !active {
			return
		}
		s.handleRestoreHealth(ev, p)
	case event.ReceiveDamage:
		if !active {
			return
		}
		s.handleReceiveDamage(ev, p)
	}
}

func (s *session) handleRegister(ev event.Event, p *domain.Player) {
	if p.Registered {
		slog.Warn("duplicate registration", "player_id", ev.PlayerID)
		s.log.Emit(ev.Time, output.MsgDisqualified(ev.PlayerID))
		p.Disqual = true
		return
	}
	p.Registered = true
	slog.Info("player registered", "player_id", ev.PlayerID)
	s.log.Emit(ev.Time, output.MsgRegistered(ev.PlayerID))
}

func (s *session) handleEnterDungeon(ev event.Event, p *domain.Player) {
	if !p.Registered || p.Disqual {
		s.log.Emit(ev.Time, output.MsgDisqualified(ev.PlayerID))
		p.Disqual = true
		return
	}
	if p.InDungeon || p.Dead || p.DungeonDone {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	if !s.schedule.IsOpen(ev.Time) {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	p.InDungeon = true
	p.EnteredAt = ev.Time
	p.CurrentFloor = 1
	p.FloorEnteredAt = ev.Time
	s.log.Emit(ev.Time, output.MsgEnteredDungeon(ev.PlayerID))
}

func (s *session) handleKillMonster(ev event.Event, p *domain.Player) {
	if p.OnBossFloor || p.CurrentFloor > s.regularFloors {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	idx := p.CurrentFloor - 1
	if p.FloorsCleared[idx] || p.MonstersLeft[idx] == 0 {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	p.MonstersLeft[idx]--
	s.log.Emit(ev.Time, output.MsgKilledMonster(ev.PlayerID))
	p.TryFloorClear(s.regularFloors, ev.Time)
}

func (s *session) handleNextFloor(ev event.Event, p *domain.Player) {
	if p.OnBossFloor || p.CurrentFloor > s.regularFloors {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	if !p.FloorsCleared[p.CurrentFloor-1] {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	p.StopFloorTimer(s.regularFloors, ev.Time)
	p.CurrentFloor++
	p.FloorEnteredAt = ev.Time
	s.log.Emit(ev.Time, output.MsgNextFloor(ev.PlayerID))
}

func (s *session) handlePrevFloor(ev event.Event, p *domain.Player) {
	if p.OnBossFloor || p.CurrentFloor <= 1 {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	p.StopFloorTimer(s.regularFloors, ev.Time)
	p.CurrentFloor--
	p.FloorEnteredAt = ev.Time
	s.log.Emit(ev.Time, output.MsgPrevFloor(ev.PlayerID))
}

func (s *session) handleEnterBoss(ev event.Event, p *domain.Player) {
	if p.OnBossFloor {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	if !p.AllRegularFloorsCleared() {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	p.StopFloorTimer(s.regularFloors, ev.Time)
	p.OnBossFloor = true
	p.BossFloorEnteredAt = ev.Time
	s.log.Emit(ev.Time, output.MsgEnteredBossFloor(ev.PlayerID))
}

func (s *session) handleKillBoss(ev event.Event, p *domain.Player) {
	if !p.OnBossFloor {
		s.log.Emit(ev.Time, output.MsgIllegalMove(ev.PlayerID, ev.ID))
		return
	}
	p.BossDefeated = true
	p.BossKillTime = ev.Time.Sub(p.BossFloorEnteredAt)
	s.log.Emit(ev.Time, output.MsgKilledBoss(ev.PlayerID))
}

func (s *session) handleLeaveDungeon(ev event.Event, p *domain.Player) {
	p.InDungeon = false
	p.DungeonDone = true
	p.LeftAt = ev.Time
	s.log.Emit(ev.Time, output.MsgLeftDungeon(ev.PlayerID))
}

func (s *session) handleCannotContinue(ev event.Event, p *domain.Player) {
	p.Disqual = true
	p.EndChallenge(ev.Time)
	s.log.Emit(ev.Time, output.MsgCannotContinue(ev.PlayerID, ev.Extra))
}

func (s *session) handleRestoreHealth(ev event.Event, p *domain.Player) {
	hp, err := strconv.Atoi(ev.Extra)
	if err != nil || hp <= 0 {
		return
	}
	p.RestoreHealth(hp)
	s.log.Emit(ev.Time, output.MsgRestoredHealth(ev.PlayerID, hp))
}

func (s *session) handleReceiveDamage(ev event.Event, p *domain.Player) {
	dmg, err := strconv.Atoi(ev.Extra)
	if err != nil || dmg <= 0 {
		return
	}
	p.ApplyDamage(dmg)
	s.log.Emit(ev.Time, output.MsgReceivedDamage(ev.PlayerID, dmg))
	if p.Health <= 0 {
		slog.Info("player died", "player_id", ev.PlayerID)
		p.MarkDead(ev.Time)
		s.log.Emit(ev.Time, output.MsgIsDead(ev.PlayerID))
	}
}
