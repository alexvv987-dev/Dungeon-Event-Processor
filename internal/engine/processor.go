package engine

import (
	"io"
	"log/slog"

	"impulse/internal/config"
	"impulse/internal/domain"
	"impulse/internal/event"
	"impulse/internal/output"
)

type session struct {
	cfg           config.Config
	regularFloors int
	schedule      domain.Schedule
	players       map[int]*domain.Player
	playerOrder   []int
	log           *output.Log
}

func newSession(cfg config.Config, schedule domain.Schedule) *session {
	return &session{
		cfg:           cfg,
		regularFloors: cfg.RegularFloors(),
		schedule:      schedule,
		players:       make(map[int]*domain.Player),
		log:           &output.Log{},
	}
}

func (s *session) fetchPlayer(id int) *domain.Player {
	if p, ok := s.players[id]; ok {
		return p
	}
	p := domain.NewPlayer(id, s.regularFloors, s.cfg)
	s.players[id] = p
	s.playerOrder = append(s.playerOrder, id)
	return p
}

func (s *session) sortedPlayers() []*domain.Player {
	out := make([]*domain.Player, 0, len(s.playerOrder))
	for _, id := range s.playerOrder {
		out = append(out, s.players[id])
	}
	return out
}

func Process(cfg config.Config, events io.Reader) (string, error) {
	if err := cfg.Validate(); err != nil {
		return "", err
	}

	schedule, err := domain.NewSchedule(cfg)
	if err != nil {
		return "", err
	}

	parsed, err := event.ReadAll(events)
	if err != nil {
		return "", err
	}

	slog.Info("starting session", "floors", cfg.Floors, "monsters", cfg.Monsters, "open_at", cfg.OpenAt, "events", len(parsed))

	sess := newSession(cfg, schedule)
	for _, ev := range parsed {
		sess.handle(ev)
	}

	domain.CloseRemainingPlayers(sess.players, schedule.CloseAt)
	output.WriteFinalReport(sess.log, sess.sortedPlayers(), sess.regularFloors)
	slog.Info("session complete", "players", len(sess.players))
	return sess.log.String(), nil
}
