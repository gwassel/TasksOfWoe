package add

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/gwassel/TasksOfWoe/internal/usecase/list/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/list/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(logger infra.Logger, db *sqlx.DB) *usecase.Usecase {
	taskRepo := repository.New(db)

	return usecase.New(logger, taskRepo)
}
