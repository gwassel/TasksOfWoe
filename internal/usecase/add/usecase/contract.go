package usecase

type TaskRepo interface {
	AddTask(userID int64, task string) (int64, error)
}
