//go:generate go tool mockgen -package $GOPACKAGE -source $GOFILE -destination contract_mock.go

package usecase

import domain "github.com/gwassel/TasksOfWoe/internal/domain/task"

type TaskRepo interface {
	ListTasks(userID int64) ([]domain.Task, error)
}
