package amqpapp

import (
	"github.com/arumandesu/uniclubs-posts-service/internal/rabbitmq"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type App struct {
	log         *slog.Logger
	amqp        Amqp
	usrService  UserService
	clubService ClubService
}

type Amqp interface {
	Consume(queue string, routingKey string, handler func(msg amqp091.Delivery) error) error
	Close() error
}

type UserService interface {
	HandleCreateUser(msg amqp091.Delivery) error
	HandleUpdateUser(msg amqp091.Delivery) error
	HandleDeleteUser(msg amqp091.Delivery) error
}

type ClubService interface {
	HandleCreateClub(msg amqp091.Delivery) error
	HandleUpdateClub(msg amqp091.Delivery) error
	HandleDeleteClub(msg amqp091.Delivery) error
}

func New(log *slog.Logger, userService UserService, clubService ClubService, amqp Amqp) *App {
	return &App{
		log:         log,
		amqp:        amqp,
		usrService:  userService,
		clubService: clubService,
	}
}

func (a *App) SetupMessageConsumers() {
	a.consumeMessages(rabbitmq.UserEventsQueue, rabbitmq.UserActivatedEventRoutingKey, a.usrService.HandleCreateUser)
	a.consumeMessages(rabbitmq.UserEventsQueue, rabbitmq.UserUpdatedEventRoutingKey, a.usrService.HandleUpdateUser)
	a.consumeMessages(rabbitmq.UserEventsQueue, rabbitmq.UserDeletedEventRoutingKey, a.usrService.HandleDeleteUser)

	a.consumeMessages(rabbitmq.ClubEventsQueue, rabbitmq.ClubCreatedEventRoutingKey, a.clubService.HandleCreateClub)
	a.consumeMessages(rabbitmq.ClubEventsQueue, rabbitmq.ClubUpdatedEventRoutingKey, a.clubService.HandleUpdateClub)
	//a.consumeMessages(rabbitmq.ClubEventsQueue, rabbitmq.ClubDeletedEventRoutingKey, a.clubService.HandleDeleteClub)
}

func (a *App) consumeMessages(queue, routingKey string, handler rabbitmq.Handler) {
	go func() {
		const op = "amqp.app.consumeMessages"
		log := a.log.With(slog.String("op", op))

		err := a.amqp.Consume(queue, routingKey, handler)
		if err != nil {
			log.Error("failed to consume ", logger.Err(err))
		}
	}()
}

func (a *App) Shutdown() {
	// Add your cleanup logic here
	// For example, close the RabbitMQ connection
	a.amqp.Close()
}
