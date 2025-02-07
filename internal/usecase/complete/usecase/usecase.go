package usecase

type Usecase struct {
	taskRepo TaskRepo
}

func New(taskRepo TaskRepo) *Usecase {
	return &Usecase{taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID, taskID int64) string {
	err := u.taskRepo.CompleteTask(userID, taskID)
	if err != nil {
		return "Failed to complete task."
	}
	return "Task marked as completed!"
}
