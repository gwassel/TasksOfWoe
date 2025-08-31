package domain

import (
	"time"
	_ "time/tzdata"
)

type Task struct {
	ID            int64   `db:"id"`
	UserTaskID    int64   `db:"user_task_id"`
	UserID        int64   `db:"user_id"`
	EncryptedTask []byte  `db:"encrypted_task"`
	Task          string  `db:"task"`
	CreatedAt     string  `db:"created_at"`
	Completed     bool    `db:"completed"`
	CompletedAt   *string `db:"completed_at"`
	InWork        bool    `db:"is_in_work"`
	TakenAt       *string `db:"taken_at"`
}

type taskStatus int

const (
	Incomplete taskStatus = iota
	Working
	Completed
)

func (t *Task) Status() taskStatus {
	status := Incomplete
	if t.Completed {
		status = Completed
	} else if t.InWork {
		status = Working
	}

	return status
}

// Converts given date in RFC3339Nano format to DateTime format Moscow time
func FormatDateForTask(date string) (string, error) {
	Moscow, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		return "", err
	}
	tm, err := time.Parse(time.RFC3339Nano, date)
	if err != nil {
		return "", err
	}

	return tm.In(Moscow).Format(time.DateTime) + " (Moscow)", nil
}

func (status taskStatus) ToString() string {
	switch status {
	case Incomplete:
		return "Incomplete"
	case Working:
		return "Working"
	case Completed:
		return "Completed"
	}
	return ""
}
