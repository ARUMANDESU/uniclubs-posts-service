package club

import (
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log"
)

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s Service) HandleCreateClub(msg amqp091.Delivery) error {
	var input struct {
		ID      int64  `json:"id"`
		Name    string `json:"name"`
		LogoUrl string `json:"logo_url"`
	}

	err := json.Unmarshal(msg.Body, &input)
	if err != nil {
		return err
	}
	log.Println(input)
	//TODO implement me
	return nil
}

func (s Service) HandleUpdateClub(msg amqp091.Delivery) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) HandleDeleteClub(msg amqp091.Delivery) error {
	//TODO implement me
	panic("implement me")
}
