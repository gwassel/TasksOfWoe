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
