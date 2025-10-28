package usecase

import (
	"context"

	"github.com/gwassel/TasksOfWoe/internal/usecase/analytics/everyday_report/domain"
)

type repository interface {
	GetFieldForReport(ctx context.Context) (domain.Metric, error)
}
