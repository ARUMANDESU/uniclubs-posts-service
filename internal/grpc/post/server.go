package postgrpc

import (
	postv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/post"
	"google.golang.org/grpc"
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
