package club

import (
	"context"
	"errors"
	"fmt"
	clubv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/club"
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
	ErrClubNotFound = errors.New("club not found")
	ErrInvalidArg   = errors.New("invalid argument")
)

type Client struct {
	clubv1.ClubClient
	log *slog.Logger
}

func (c *Client) IsClubMember(ctx context.Context, userId, clubId int64) (bool, error) {
	const op = "client.club.isClubMember"
	log := c.log.With(slog.String("op", op))

	res, err := c.ClubClient.GetJoinStatus(ctx, &clubv1.GetJoinStatusRequest{
		ClubId: clubId,
		UserId: userId,
	})
	if err != nil {
		switch {
		case codes.InvalidArgument == status.Code(err):
			return false, ErrInvalidArg
		case codes.NotFound == status.Code(err):
			return false, ErrClubNotFound
		default:
			log.Error("internal", logger.Err(err))
			return false, err
		}
	}

	if res.GetStatus() == clubv1.JoinStatus_MEMBER {
		return true, nil
	}

	return false, nil

}

func (c *Client) IsBanned(ctx context.Context, userId, clubId int64) (bool, error) {
	const op = "client.club.isBanned"
	log := c.log.With(slog.String("op", op))

	res, err := c.ClubClient.GetJoinStatus(ctx, &clubv1.GetJoinStatusRequest{
		ClubId: clubId,
		UserId: userId,
	})
	if err != nil {
		switch {
		case codes.InvalidArgument == status.Code(err):
			return false, ErrInvalidArg
		case codes.NotFound == status.Code(err):
			return false, ErrClubNotFound
		default:
			log.Error("internal", logger.Err(err))
			return false, err
		}
	}

	if res.GetStatus() == clubv1.JoinStatus_BANNED {
		return true, nil
	}

	return false, nil

}

func (c *Client) HasPermission(ctx context.Context, userId, clubId int64, permission clubv1.Permission) (bool, error) {
	const op = "client.club.hasPermission"
	log := c.log.With(slog.String("op", op))

	res, err := c.ClubClient.HavePermissionTo(ctx, &clubv1.HavePermissionToRequest{
		ClubId:     clubId,
		UserId:     userId,
		Permission: permission,
	})
	if err != nil {
		switch {
		case codes.InvalidArgument == status.Code(err):
			return false, ErrInvalidArg
		case codes.NotFound == status.Code(err):
			return false, ErrClubNotFound
		default:
			log.Error("internal", logger.Err(err))
			return false, err
		}
	}

	return res.GetHasPermission(), nil
}

func (c *Client) GetClubById(ctx context.Context, clubId int64) (*domain.Club, error) {
	const op = "client.club.getClubById"
	log := c.log.With(slog.String("op", op))

	res, err := c.ClubClient.GetClub(ctx, &clubv1.GetClubRequest{ClubId: clubId})
	if err != nil {
		switch {
		case codes.InvalidArgument == status.Code(err):
			return nil, ErrInvalidArg
		case codes.NotFound == status.Code(err):
			return nil, ErrClubNotFound
		default:
			log.Error("internal", logger.Err(err))
			return nil, err
		}
	}

	return &domain.Club{
		ID:      res.GetClubId(),
		Name:    res.GetName(),
		LogoURL: res.GetLogoUrl(),
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
		ClubClient: clubv1.NewClubClient(cc),
		log:        log,
	}, nil
}

// InterceptorLogger adapts slog logger to interceptor logger
func InterceptorLogger(l *slog.Logger) grpclog.Logger {
	return grpclog.LoggerFunc(func(ctx context.Context, lvl grpclog.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
