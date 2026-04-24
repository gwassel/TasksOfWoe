//go:generate go tool mockgen -package $GOPACKAGE -source $GOFILE -destination contract_mock.go

package usecase

type TaskRepo interface {
	AddTask(userID int64, task string) (int64, error)
}
