package usecase

import (
	"context"
	"time"

	"github.com/gwassel/TasksOfWoe/internal/infra"
	"github.com/gwassel/TasksOfWoe/internal/usecase/analytics/everyday_report/domain"
	"github.com/pkg/errors"
)

type usecase struct {
	logger infra.Logger
	repos  []repository
}

func New(logger infra.Logger, repos []repository) *usecase {
	return &usecase{
		logger: logger,
		repos:  repos,
	}
}

func (u *usecase) Handle(ctx context.Context) (domain.EverydayReportResponse, error) {
	result := domain.EverydayReportResponse{
		Date:                 time.Now().Add(-24 * time.Hour),
		UserIDsToSendMetrics: []int64{351083864},
	}

	for _, r := range u.repos {
		metric, err := r.GetFieldForReport(ctx)
		if err != nil {
			return result, errors.Wrap(err, "unable to calculate metrics")
		}

		result.Metrics = append(result.Metrics, metric)
	}

	return result, nil
}
