package output

import (
	"fmt"
	"strings"
	"time"

	"impulse/internal/domain"
	"impulse/internal/timeclock"
)

type Log struct {
	buf strings.Builder
}

func (l *Log) Emit(t time.Time, msg string) {
	fmt.Fprintf(&l.buf, "[%s] %s\n", timeclock.FormatTime(t), msg)
}

func (l *Log) WriteLine(line string) {
	fmt.Fprintln(&l.buf, line)
}

func (l *Log) String() string {
	return l.buf.String()
}

func WriteFinalReport(log *Log, players []*domain.Player, regularFloors int) {
	log.WriteLine("Final report:")
	for _, p := range players {
		log.WriteLine(fmt.Sprintf("[%s] %d [%s, %s, %s] HP:%d",
			p.FinalState(), p.ID,
			timeclock.FormatDuration(p.DungeonDuration()),
			timeclock.FormatDuration(p.AvgFloorTime(regularFloors)),
			timeclock.FormatDuration(p.BossKillTime),
			p.Health,
		))
	}
}

func MsgRegistered(id int) string {
	return fmt.Sprintf("Player [%d] registered", id)
}

func MsgDisqualified(id int) string {
	return fmt.Sprintf("Player [%d] disqualified", id)
}

func MsgEnteredDungeon(id int) string {
	return fmt.Sprintf("Player [%d] entered the dungeon", id)
}

func MsgKilledMonster(id int) string {
	return fmt.Sprintf("Player [%d] killed the monster", id)
}

func MsgNextFloor(id int) string {
	return fmt.Sprintf("Player [%d] went to the next floor", id)
}

func MsgPrevFloor(id int) string {
	return fmt.Sprintf("Player [%d] went to the previous floor", id)
}

func MsgEnteredBossFloor(id int) string {
	return fmt.Sprintf("Player [%d] entered the boss's floor", id)
}

func MsgKilledBoss(id int) string {
	return fmt.Sprintf("Player [%d] killed the boss", id)
}

func MsgLeftDungeon(id int) string {
	return fmt.Sprintf("Player [%d] left the dungeon", id)
}

func MsgCannotContinue(id int, reason string) string {
	return fmt.Sprintf("Player [%d] cannot continue due to [%s]", id, reason)
}

func MsgRestoredHealth(id, amount int) string {
	return fmt.Sprintf("Player [%d] has restored [%d] of health", id, amount)
}

func MsgReceivedDamage(id, dmg int) string {
	return fmt.Sprintf("Player [%d] recieved [%d] of damage", id, dmg)
}

func MsgIsDead(id int) string {
	return fmt.Sprintf("Player [%d] is dead", id)
}

func MsgIllegalMove(id, eventID int) string {
	return fmt.Sprintf("Player [%d] makes imposible move [%d]", id, eventID)
}
