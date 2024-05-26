package userclient

import (
	"context"
	"errors"
	"fmt"
	userv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/user"
	"github.com/arumandesu/uniclubs-posts-service/internal/domain"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidArg   = errors.New("invalid argument")
)

type Client struct {
	userv1.UserClient
	log *slog.Logger
}

func (c *Client) GetUserById(ctx context.Context, id int64) (*domain.User, error) {
	const op = "client.user.getUserByID"
	log := c.log.With(slog.String("op", op))

	user, err := c.UserClient.GetUser(ctx, &userv1.GetUserRequest{UserId: id})
	if err != nil {
		switch {
		case status.Code(err) == codes.InvalidArgument:
			return nil, ErrInvalidArg
		case status.Code(err) == codes.NotFound:
			return nil, ErrUserNotFound
		default:
			log.Error("internal", logger.Err(err))
			return nil, err
		}
	}

	return &domain.User{
		ID:        user.GetUserId(),
		FirstName: user.GetFirstName(),
		LastName:  user.GetLastName(),
		Barcode:   user.GetBarcode(),
		AvatarURL: user.GetAvatarUrl(),
	}, nil
}

func New(
	log *slog.Logger,
	addr string,
	timeout time.Duration,
	retriesCount int,
) (*Client, error) {
	const op = "grpc.New"

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.NotFound, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	logOpts := []grpclog.Option{
		grpclog.WithLogOnEvents(grpclog.StartCall, grpclog.FinishCall),
	}

	cc, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // in the future, we can use tls/ssl cert if we want
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(InterceptorLogger(log), logOpts...),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Client{
		UserClient: userv1.NewUserClient(cc),
		log:        log,
	}, nil
}

// InterceptorLogger adapts slog logger to interceptor logger
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
