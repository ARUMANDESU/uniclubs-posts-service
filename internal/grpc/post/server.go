package postgrpc

import (
	"errors"
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	postservice "github.com/arumandesu/uniclubs-posts-service/internal/services/post"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type serverApi struct {
	postv1.PostServer
	Services
}

func (s serverApi) mustEmbedUnimplementedPostServer() {
	//TODO implement me
	panic("implement me")
}

type Services struct {
	management ManagementService
	info       InfoService
}

func Register(
	gRPC *grpc.Server,
	services Services,
) {
	postv1.RegisterPostServer(gRPC, &serverApi{Services: services})
}

func NewServices(
	management ManagementService,
	info InfoService,
) Services {
	return Services{
		management: management,
		info:       info,
	}
}

func handleServiceError(err error) error {
	switch {
	case errors.Is(err, postservice.ErrPostNotFound), errors.Is(err, postservice.ErrClubNotFound):
		return status.Error(codes.NotFound, err.Error())
	case errors.Is(err, postservice.ErrPermissionDenied):
		return status.Error(codes.PermissionDenied, err.Error())
	case errors.Is(err, postservice.ErrInvalidArg):
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, "internal error")
	}
}
