package encode

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/gwassel/TasksOfWoe/internal/usecase/encode/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/encode/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(logger infra.Logger, db *sqlx.DB, key string) *usecase.Usecase {
	taskRepo := repository.New(db)

	return usecase.New(logger, taskRepo, key)
}
