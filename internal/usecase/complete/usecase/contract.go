package usecase

type TaskRepo interface {
	CompleteTask(userID int64, taskID int64) error
}
