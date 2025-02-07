package complete

import (
	"github.com/gwassel/TasksOfWoe/internal/usecase/complete/repository"
	"github.com/gwassel/TasksOfWoe/internal/usecase/complete/usecase"
	"github.com/jmoiron/sqlx"
)

func NewUsecase(db *sqlx.DB) *usecase.Usecase {
	taskRepo := repository.New(db)

	return usecase.New(taskRepo)
}
