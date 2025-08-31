package add

import (
	"github.com/gwassel/TasksOfWoe/internal/usecase/listall/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/listall/service/encoder"
	"github.com/gwassel/TasksOfWoe/internal/usecase/listall/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(db *sqlx.DB, taskEncoder encoder.Encoder) *usecase.Usecase {
	taskRepo := repository.New(db)
	encoderTaskRepo := encoder.New(taskRepo, taskEncoder)

	return usecase.New(encoderTaskRepo)
}
