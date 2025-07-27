package usecase

type TaskRepo interface {
	ToggleInWork(userID int64, taskID int64) error
}
