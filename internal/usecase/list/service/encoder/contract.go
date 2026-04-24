package encoder

import (
	"context"

	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
)

type TasksRepository interface {
	ListTasks(ctx context.Context, userID int64) ([]domain.Task, error)
}

type Encoder interface {
	Decode(ciphertext []byte) (string, error)
}
