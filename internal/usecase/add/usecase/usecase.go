package usecase

type Usecase struct {
	taskRepo TaskRepo
}

func New(taskRepo TaskRepo) *Usecase {
	return &Usecase{taskRepo: taskRepo}
}

func (u *Usecase) Handle(userID int64, task string) string {
	err := u.taskRepo.AddTask(userID, task)
	if err != nil {
		return "Failed to add task."
	}
	return "Task added."
}
