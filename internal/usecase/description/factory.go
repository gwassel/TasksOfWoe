package description

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/gwassel/TasksOfWoe/internal/usecase/description/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/description/service/encoder"
	"github.com/gwassel/TasksOfWoe/internal/usecase/description/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(logger infra.Logger, db *sqlx.DB, taskEncoder encoder.Encoder) *usecase.Usecase {
	taskRepo := repository.New(db)
	encoderTaskRepo := encoder.New(taskRepo, taskEncoder)

	return usecase.New(logger, encoderTaskRepo)
}
