package user

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

func (s Service) HandleCreateUser(msg amqp091.Delivery) error {
	var input struct {
		ID        int64  `json:"id"`
		Barcode   string `json:"barcode"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		AvatarURL string `json:"avatar_url"`
	}

	err := json.Unmarshal(msg.Body, &input)
	if err != nil {
		return err
	}
	log.Println(input)
	//TODO implement me
	return nil
}

func (s Service) HandleUpdateUser(msg amqp091.Delivery) error {
	//TODO implement me
	panic("implement me")
}

func (s Service) HandleDeleteUser(msg amqp091.Delivery) error {
	//TODO implement me
	panic("implement me")
}
