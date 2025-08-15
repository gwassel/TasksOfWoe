package usecase

type TaskRepo interface {
	ToggleInWork(userID int64, taskIDs []int64, value bool) error
}
