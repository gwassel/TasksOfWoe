package usecase

type TaskRepo interface {
	AddTask(userID int64, task []byte) (int64, error)
}
