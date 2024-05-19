package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arumandesu/uniclubs-posts-service/internal/config"
	"github.com/arumandesu/uniclubs-posts-service/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
)

const (
	ClubExchangeName           = "club-exchange"
	UserExchangeName           = "user-exchange"
	UserEventsQueue            = "user-events-posts-queue"
	ClubEventsQueue            = "club-events-posts-queue"
	UserUpdatedEventRoutingKey = "user.event.updated"
	ClubUpdatedEventRoutingKey = "club.event.updated"
)

type Handler func(msg amqp.Delivery) error

type Rabbitmq struct {
	conn *amqp.Connection
	ch   *amqp.Channel
	cfg  config.Rabbitmq
	log  *slog.Logger
}

func New(cfg config.Rabbitmq, log *slog.Logger) (*Rabbitmq, error) {
	const op = "rabbitmq.new"

	connString := fmt.Sprintf("amqp://%v:%v@%v:%v/", cfg.User, cfg.Password, cfg.Host, cfg.Port)
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to connect to amqp server: %w", op, err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("%s: failed to open a channel: %w", op, err)
	}

	err = declareExchanges(ch)
	if err != nil {
		log.Error("failed to declare exchanges", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = declareQueues(ch)
	if err != nil {
		log.Error("failed to declare queues", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = bindQueues(ch)
	if err != nil {
		log.Error("failed to bind queues", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Rabbitmq{
		conn: conn,
		ch:   ch,
		cfg:  cfg,
		log:  log,
	}, nil
}

func (r *Rabbitmq) Consume(queue, routingKey string, handler func(msg amqp.Delivery) error) error {
	const op = "rabbitmq.consume"
	log := r.log.With(
		slog.String("op", op),
		slog.With(
			slog.String("queue", queue),
			slog.String("routing_key", routingKey),
		),
	)

	err := r.ch.Qos(
		1,
		0,
		false,
	)
	if err != nil {
		log.Error("failed to set Qos", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	msgs, err := r.ch.Consume(
		queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error("failed to register as consumer", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Debug("routing key", slog.String("key", d.RoutingKey))
			if d.RoutingKey != routingKey {
				err := d.Reject(true)
				if err != nil {
					log.Warn("failed to negatively acknowledge", logger.Err(err))
				}
				continue
			}

			err = handler(d)
			if err != nil {
				log.Warn("failed to handle message", logger.Err(err))
				continue
			}
			err = d.Ack(false)
			if err != nil {
				log.Warn("failed to send an acknowledgement", logger.Err(err))
			}

		}
	}()

	<-forever

	return nil
}

func (r *Rabbitmq) Publish(ctx context.Context, exchangeName string, routingKey string, msg any) error {
	const op = "rabbitmq.publish"

	bytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = r.ch.PublishWithContext(
		ctx,
		exchangeName,
		routingKey,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         bytes,
		})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Rabbitmq) Close() error {
	const op = "rabbitmq.close"
	log := r.log.With(slog.String("op", op))

	err := r.ch.Close()
	if err != nil {
		log.Error("failed to close channel", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	err = r.conn.Close()
	if err != nil {
		log.Error("failed to close connection", logger.Err(err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func declareExchanges(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(
		ClubExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.ExchangeDeclare(
		UserExchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func declareQueues(ch *amqp.Channel) error {
	_, err := ch.QueueDeclare(
		UserEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	_, err = ch.QueueDeclare(
		ClubEventsQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}

func bindQueues(ch *amqp.Channel) error {
	err := ch.QueueBind(
		UserEventsQueue,
		UserUpdatedEventRoutingKey,
		UserExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.QueueBind(
		ClubEventsQueue,
		ClubUpdatedEventRoutingKey,
		ClubExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return nil
}
