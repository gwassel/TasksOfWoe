package usecase

import (
	"context"
	"testing"

	domain "github.com/gwassel/TasksOfWoe/internal/domain/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListUsecase_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().ListTasks(ctx, int64(1)).Return([]domain.Task{{}}, nil)

	usecase := New(nil, mockRepo)
	tasks, err := usecase.Handle(ctx, int64(1))
	require.NoError(t, err)
	require.Len(t, tasks, 1)
}

func TestListUsecase_Handle_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().ListTasks(ctx, int64(1)).Return([]domain.Task{}, assert.AnError)

	usecase := New(nil, mockRepo)
	tasks, err := usecase.Handle(ctx, int64(1))
	require.Error(t, err)
	require.Len(t, tasks, 0)
}
