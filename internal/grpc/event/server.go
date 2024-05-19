package event

import (
	eventv1 "github.com/ARUMANDESU/uniclubs-protos/gen/go/posts/event"
	"google.golang.org/grpc"
)

type serverApi struct {
	eventv1.UnimplementedEventServer
	Services
}

type Services struct {
	management   ManagementService
	organizer    OrganizerService
	collaborator CollaboratorService
	info         InfoService
}

func Register(
	gRPC *grpc.Server,
	services Services,
) {
	eventv1.RegisterEventServer(gRPC, &serverApi{Services: services})
}

func NewServices(
	management ManagementService,
	organizer OrganizerService,
	collaborator CollaboratorService,
	info InfoService,
) Services {
	return Services{
		management:   management,
		organizer:    organizer,
		collaborator: collaborator,
		info:         info,
	}
}
