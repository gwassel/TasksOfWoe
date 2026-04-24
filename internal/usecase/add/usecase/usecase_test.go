package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAddUsecase_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().AddTask(ctx, int64(1), "test task").Return(int64(101), nil)

	usecase := New(mockRepo)
	userTaskID, err := usecase.Handle(ctx, int64(1), "test task")
	require.NoError(t, err)
	require.Equal(t, int64(101), userTaskID)
}

func TestAddUsecase_Handle_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().AddTask(ctx, int64(1), "test task").Return(int64(0), assert.AnError)

	usecase := New(mockRepo)
	userTaskID, err := usecase.Handle(ctx, int64(1), "test task")
	require.Error(t, err)
	require.Equal(t, int64(0), userTaskID)
}
