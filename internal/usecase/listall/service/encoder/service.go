package encoder

import (
	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/pkg/errors"
)

type EncoderService struct {
	taskRepo TasksRepository
	encoder  Encoder
}

func New(
	taskRepo TasksRepository,
	encoder Encoder,
) *EncoderService {
	return &EncoderService{
		taskRepo: taskRepo,
		encoder:  encoder,
	}
}

func (es *EncoderService) ListAllTasks(userID int64) ([]domain.Task, error) {
	tasks, err := es.taskRepo.ListAllTasks(userID)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to list tasks")
	}

	for i, task := range tasks {
		tasks[i].Task, err = es.encoder.Decode(task.EncryptedTask)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to decodee tasks")
		}
	}

	return tasks, nil
}
