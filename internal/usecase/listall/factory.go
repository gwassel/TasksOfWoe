package add

import (
	"github.com/gwassel/TasksOfWoe/internal/usecase/listall/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/listall/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(db *sqlx.DB) *usecase.Usecase {
	taskRepo := repository.New(db)

	return usecase.New(taskRepo)
}
