package usecase

type TaskRepo interface {
	TakeTask(userID int64, taskIDs []int64) error
}
