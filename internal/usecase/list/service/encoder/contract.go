package encoder

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type TasksRepository interface {
	ListTasks(userID int64) ([]domain.Task, error)
}

type Encoder interface {
	Decode(ciphertext []byte) (string, error)
}
