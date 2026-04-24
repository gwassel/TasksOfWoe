//go:generate go tool mockgen -package $GOPACKAGE -source $GOFILE -destination contract_mock.go

package usecase

type TaskRepo interface {
	UntakeTask(userID int64, taskIDs []int64) error
}
