//go:generate go tool mockgen -package $GOPACKAGE -source $GOFILE -destination contract_mock.go

package usecase

import (
	"context"

	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
)

type TaskRepo interface {
	TaskDescription(ctx context.Context, userID int64, userTaskIDs []int64) ([]domain.Task, error)
}
