package encoder

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type EncoderService struct {
	tasksRepo TasksRepository
	encoder   Encoder
}

func New(
	tasksRepo TasksRepository,
	encoder Encoder,
) *EncoderService {
	return &EncoderService{
		tasksRepo: tasksRepo,
		encoder:   encoder,
	}
}

func (es *EncoderService) ListTasks(userID int64) ([]domain.Task, error) {
	tasks, err := es.ListTasks(userID)
	if err != nil {
		return nil, err // wrap
	}

	for i, task := range tasks {
		tasks[i].Text, err = es.encoder.Decode(task.Task)
		if err != nil {
			return nil, err
		}
	}

	return tasks, nil
}
