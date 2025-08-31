package encoder

type TasksRepository interface {
	AddTask(userID int64, task []byte) (int64, error)
}

type Encoder interface {
	Encode(plaintext string) ([]byte, error)
}
