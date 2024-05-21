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

type suite struct {
	UserService *Service
	*mockStorage
}

type mockStorage struct {
	mock.Mock
}

func (m *mockStorage) UpdateUser(ctx context.Context, user *domain.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func setupSuite(t *testing.T) *suite {
	t.Helper()

	mockStorage := &mockStorage{}
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	userService := New(log, mockStorage)

	return &suite{
		UserService: userService,
		mockStorage: mockStorage,
	}
}

func TestService_HandleUpdateUser(t *testing.T) {
	suite := setupSuite(t)
	defer suite.AssertExpectations(t)

	t.Run("success", func(t *testing.T) {

		user := &domain.User{
			ID:        1,
			FirstName: "John Doe",
		}
		newDelivery := amqp.Delivery{Body: []byte(`{"id":1,"first_name":"John Doe"}`)}

		suite.mockStorage.On("UpdateUser", mock.Anything, user).Return(nil)

		err := suite.UserService.HandleUpdateUser(newDelivery)
		assert.NoError(t, err)
	})

	t.Run("user not found", func(t *testing.T) {

		user := &domain.User{
			ID:        3,
			FirstName: "John Doe",
		}
		newDelivery := amqp.Delivery{Body: []byte(`{"id":3,"first_name":"John Doe"}`)}

		suite.mockStorage.On("UpdateUser", mock.Anything, user).Return(storage.ErrUserNotExists)

		err := suite.UserService.HandleUpdateUser(newDelivery)
		t.Logf("err: %v, want: %v", err, ErrUserNotExist)
		assert.True(t, errors.Is(err, ErrUserNotExist))
	})
}
