//go:generate go tool mockgen -package $GOPACKAGE -source $GOFILE -destination contract_mock.go

package usecase

import "context"

type TaskRepo interface {
	AddTask(ctx context.Context, userID int64, task string) (int64, error)
}
