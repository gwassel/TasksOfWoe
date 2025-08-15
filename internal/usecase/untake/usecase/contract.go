package usecase

type TaskRepo interface {
	UntakeTask(userID int64, taskIDs []int64) error
}
