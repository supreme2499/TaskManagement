package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"

	"Tasks/internal/lib/logger/handler/slogdiscard"
	"Tasks/internal/model"
	mockery "Tasks/internal/service/mocks"
)

type mocks struct {
	repositoryStorage *mockery.StorageRepository
	repositoryCache   *mockery.CacheRepository
	broker            *mockery.Broker
}

func TestService_CreateTask(t *testing.T) {
	location, _ := time.LoadLocation("Europe/Moscow")

	tests := []struct {
		name     string
		input    model.Task
		expected int
		wantErr  bool
		mock     func() mocks
	}{
		{
			name:     "positive base test",
			input:    model.Task{NameTask: "task123", Description: "opisanie", Deadline: time.Date(2025, 3, 30, 10, 0, 0, 0, location)},
			expected: 0,
			mock: func() mocks {
				storageMock := mockery.NewStorageRepository(t)
				storageMock.On("CreateNewTask", mock.Anything, mock.Anything).Return(0, nil)

				cacheMock := mockery.NewCacheRepository(t)
				cacheMock.On("InsertingCache", mock.Anything, mock.Anything).Return(nil)

				return mocks{repositoryStorage: storageMock, repositoryCache: cacheMock}
			},
		},
		{name: "negative test 1", input: model.Task{
			NameTask:    "task123",
			Description: "opisanie",
			Deadline:    time.Date(2024, 3, 30, 10, 0, 0, 0, location)},
			expected: -1,
			wantErr:  true,
			mock: func() mocks {
				storageMock := mockery.NewStorageRepository(t)
				cacheMock := mockery.NewCacheRepository(t)

				return mocks{repositoryStorage: storageMock, repositoryCache: cacheMock}
			},
		},
		{name: "negative test 2", input: model.Task{
			NameTask:    "task123",
			Description: "opisanie",
			Deadline:    time.Date(2025, 2, 30, 10, 0, 0, 0, location)},
			expected: -1,
			wantErr:  true,
			mock: func() mocks {
				storageMock := mockery.NewStorageRepository(t)
				storageMock.On("CreateNewTask", mock.Anything, mock.Anything).Return(0, errors.New("incorrect date"))
				cacheMock := mockery.NewCacheRepository(t)

				return mocks{repositoryStorage: storageMock, repositoryCache: cacheMock}
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			m := tt.mock()

			ct := Service{
				log:   slogdiscard.NewDiscardLogger(),
				repo:  m.repositoryStorage,
				cache: m.repositoryCache,
			}

			_, err := ct.CreateTask(context.Background(), tt.input)
			if tt.wantErr != (err != nil) {
				t.Errorf("unexpected error: %v", err)
				return
			}

		})
	}
}
