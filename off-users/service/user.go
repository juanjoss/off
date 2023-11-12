package service

import (
	"encoding/json"
	"log"

	"github.com/juanjoss/off-users-service/ports"
)

type UserService struct {
	repo   ports.UserRepository
	pubsub ports.PubSub
}

func NewUserService(repo ports.UserRepository, pubsub ports.PubSub) *UserService {
	return &UserService{
		repo:   repo,
		pubsub: pubsub,
	}
}

func (us *UserService) Register(request ports.RegisterRequest) error {
	err := us.repo.Register(request)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) SubscribeOrdersNew() {
	_, err := us.pubsub.Subscribe("orders.random", func(data any) {
		ssdId, err := us.repo.RandomSSDId()
		if err != nil {
			log.Printf("unable to get random SSD: %v", err)
		}

		response := map[string]int{"ssd_id": ssdId}

		err = us.pubsub.Publish("orders.new", response)
		if err != nil {
			log.Printf("error while publishing to orders.random: %v", err)
		}
	})
	if err != nil {
		log.Printf("unable to subscribe to orders.random: %v", err)
	}
}

func (us *UserService) SubscribeOrdersCompleted() {
	_, err := us.pubsub.Subscribe("orders.completed", func(data any) {
		var order ports.AddProductToSSD

		err := json.Unmarshal(data.([]byte), &order)
		if err != nil {
			log.Printf("error during orders.completed: %v", err)
		}

		// add product to SSD
		err = us.repo.AddProductToSSD(order.SsdId, order.Barcode, order.Quantity)
		if err != nil {
			log.Printf("unable to add product to ssd: %v", err)
		}
	})
	if err != nil {
		log.Printf("unable to subscribe to orders.completed: %v", err)
	}
}
