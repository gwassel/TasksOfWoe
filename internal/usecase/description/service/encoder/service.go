package encoder

import (
	"context"

	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/pkg/errors"
)

type encoderService struct {
	taskRepository TaskRepository
	encoder        Encoder
}

func New(
	taskRepository TaskRepository,
	encoder Encoder,
) *encoderService {
	return &encoderService{taskRepository: taskRepository, encoder: encoder}
}

func (es *encoderService) TaskDescription(
	ctx context.Context,
	userID int64,
	userTaskIDs []int64,
) ([]domain.Task, error) {
	tasks, err := es.taskRepository.TaskDescription(ctx, userID, userTaskIDs)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to get task description")
	}

	for i, task := range tasks {
		plaintext, err := es.encoder.Decode(task.EncryptedTask)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to decrypt task(s)")
		}
		tasks[i].Task = plaintext
	}

	return tasks, nil
}
