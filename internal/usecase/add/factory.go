package add

import (
	"github.com/gwassel/TasksOfWoe/internal/usecase/add/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/add/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(db *sqlx.DB) *usecase.Usecase {
	taskRepo := repository.New(db)

	return usecase.New(taskRepo)
}
