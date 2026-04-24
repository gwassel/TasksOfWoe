package usecase

import (
	"context"
	"testing"

	"github.com/gwassel/TasksOfWoe/internal/usecase/analytics/everyday_report/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAnalyticsEverydayReportUsecase_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockrepository(ctrl)
	mockRepo.EXPECT().GetFieldForReport(gomock.Any()).Return(domain.Metric{}, nil)

	usecase := New(nil, []repository{mockRepo})
	result, err := usecase.Handle(context.Background())
	require.NoError(t, err)
	require.NotNil(t, result)
}

func TestAnalyticsEverydayReportUsecase_Handle_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockrepository(ctrl)
	mockRepo.EXPECT().GetFieldForReport(gomock.Any()).Return(domain.Metric{}, assert.AnError)

	usecase := New(nil, []repository{mockRepo})
	_, err := usecase.Handle(context.Background())
	require.Error(t, err)
}
