package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestCompleteUsecase_Handle(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().CompleteTask(int64(1), []int64{101}).Return(nil)

	usecase := New(nil, mockRepo)
	err := usecase.Handle(int64(1), []int64{101})
	require.NoError(t, err)
}

func TestCompleteUsecase_Handle_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := NewMockTaskRepo(ctrl)
	mockRepo.EXPECT().CompleteTask(int64(1), []int64{101}).Return(assert.AnError)

	usecase := New(nil, mockRepo)
	err := usecase.Handle(int64(1), []int64{101})
	require.Error(t, err)
}
