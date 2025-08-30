package add

import (
	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/gwassel/TasksOfWoe/internal/usecase/list/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/list/service/encoder"
	"github.com/gwassel/TasksOfWoe/internal/usecase/list/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(logger infra.Logger, db *sqlx.DB, taskDecoder encoder.Encoder) *usecase.Usecase {
	taskRepo := repository.New(db)
	decoderTaskRepo := encoder.New(taskRepo, taskDecoder)

	return usecase.New(logger, decoderTaskRepo)
}
