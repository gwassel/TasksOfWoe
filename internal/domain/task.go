package domain

type Task struct {
	ID          int64   `db:"id"`
	UserTaskID  int64   `db:"user_task_id"`
	UserID      int64   `db:"user_id"`
	Task        string  `db:"task"`
	CreatedAt   string  `db:"created_at"`
	Completed   bool    `db:"completed"`
	CompletedAt *string `db:"completed_at"`
	InWork      bool    `db:"is_in_work"`
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

func ToString(status taskStatus) string {
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
