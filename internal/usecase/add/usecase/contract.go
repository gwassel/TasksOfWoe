package usecase

type TaskRepo interface {
	AddTask(userID int64, task string) error
}
