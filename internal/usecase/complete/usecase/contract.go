package usecase

type TaskRepo interface {
	CompleteTask(userID int64, taskIDs []int64) error
}
