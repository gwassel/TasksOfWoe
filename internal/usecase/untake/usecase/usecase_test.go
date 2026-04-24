package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestUntakeUsecase_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().UntakeTask(ctx, int64(1), []int64{101}).Return(nil)

	usecase := New(nil, mockRepo)
	err := usecase.Handle(ctx, int64(1), []int64{101})
	require.NoError(t, err)
}

func TestUntakeUsecase_Handle_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	ctx := context.Background()
	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().UntakeTask(ctx, int64(1), []int64{101}).Return(assert.AnError)

	usecase := New(nil, mockRepo)
	err := usecase.Handle(ctx, int64(1), []int64{101})
	require.Error(t, err)
}
