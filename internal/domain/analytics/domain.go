package analytics

import (
	"time"
)

type Event struct {
	ID        int64
	TgUserID  int64
	EventName string
	Timestamp time.Time
}

func NewEvent(TgUserID int64, EventName string, Timestamp time.Time) Event {
	return Event{
		TgUserID:  TgUserID,
		EventName: EventName,
		Timestamp: Timestamp,
	}
}
