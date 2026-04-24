//go:generate go tool mockgen -package $GOPACKAGE -source $GOFILE -destination contract_mock.go

package usecase

import "context"

type TaskRepo interface {
	TakeTask(ctx context.Context, userID int64, taskIDs []int64) error
}
