package encoder

import "context"

type TasksRepository interface {
	AddTask(ctx context.Context, userID int64, task []byte) (int64, error)
}

type Encoder interface {
	Encode(plaintext string) ([]byte, error)
}
