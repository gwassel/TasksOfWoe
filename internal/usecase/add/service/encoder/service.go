package encoder

import (
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

func (es *EncoderService) AddTask(userID int64, task string) (int64, error) {
	ciphertext, err := es.encoder.Encode(task)
	if err != nil {
		return 0, errors.Wrap(err, "Failed to encode task")
	}

	userTaskID, err := es.taskRepo.AddTask(userID, ciphertext)
	if err != nil {
		return 0, errors.Wrap(err, "Faield to add task")
	}
	return userTaskID, nil
}
