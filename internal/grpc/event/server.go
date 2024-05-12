package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"google.golang.org/grpc"
)

type serverApi struct {
	eventv1.UnimplementedEventServer
	management ManagementService
}

func Register(gRPC *grpc.Server, management ManagementService) {
	eventv1.RegisterEventServer(gRPC, &serverApi{management: management})
}
