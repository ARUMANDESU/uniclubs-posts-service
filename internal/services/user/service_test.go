package userservice

import (
	"context"
	"errors"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/internal/storage"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"log/slog"
	"testing"
)

type Suite struct {
	UserService *Service
	*MockStorage
}

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func Setup(t *testing.T) *Suite {
	t.Helper()

	mockStorage := &MockStorage{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	userService := New(log, mockStorage)

	return &Suite{
		UserService: userService,
		MockStorage: mockStorage,
	}
}

func TestService_HandleUpdateUser(t *testing.T) {
	suite := Setup(t)
	defer suite.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {

		user := &domain.User{
			ID:        1,
			FirstName: "John Doe",
		}
		newDelivery := amqp.Delivery{Body: []byte(`{"id":1,"first_name":"John Doe"}`)}

		suite.MockStorage.On("UpdateUser", mock.Anything, user).Return(nil)

		err := suite.UserService.HandleUpdateUser(newDelivery)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {

		user := &domain.User{
			ID:        3,
			FirstName: "John Doe",
		}
		newDelivery := amqp.Delivery{Body: []byte(`{"id":3,"first_name":"John Doe"}`)}

		suite.MockStorage.On("UpdateUser", mock.Anything, user).Return(storage.ErrUserNotExists)

		err := suite.UserService.HandleUpdateUser(newDelivery)
		t.Logf("err: %v, want: %v", err, ErrUserNotExist)
		assert.True(t, errors.Is(err, ErrUserNotExist))
	})
}
