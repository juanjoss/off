package service

import (
	"encoding/json"
	"log"
	"time"

	"github.com/juanjoss/off-orders-service/model"
	"github.com/juanjoss/off-orders-service/ports"
)

type OrderService struct {
	repo   ports.OrderRepository
	pubsub ports.PubSub
}

func NewOrdersService(repo ports.OrderRepository, pubsub ports.PubSub) *OrderService {
	return &OrderService{
		repo:   repo,
		pubsub: pubsub,
	}
}

func (os *OrderService) CreateProductOrder(request ports.CreateProductOrderRequest) error {
	order := model.Order{
		SsdId:     request.SsdId,
		ProductId: request.Barcode,
		Timestamp: time.Now(),
		Quantity:  request.Quantity,
		Status:    "pending",
	}

	id, err := os.repo.CreateProductOrder(order)
	if err != nil {
		return err
	}

	order.Id = id

	if err := os.pubsub.Publish("orders.pending", order); err != nil {
		log.Printf("unable to publish to orders.pending: %v", err)
	}

	return nil
}

func (os *OrderService) SubscribeOrdersNew() {
	_, err := os.pubsub.Subscribe("orders.random", func(data any) {
		barcode, err := os.repo.GetRandomProductBarcode()
		if err != nil {
			log.Printf("unable to get random product barcode: %v", err)
		}

		response := map[string]string{"barcode": barcode}

		err = os.pubsub.Publish("orders.new", response)
		if err != nil {
			log.Printf("error while publishing to orders.random: %v", err)
		}
	})
	if err != nil {
		log.Printf("unable to get message from orders.random: %v", err)
	}
}

func (os *OrderService) SubscribeOrdersShipped() {
	_, err := os.pubsub.Subscribe("orders.shipped", func(data any) {
		var order model.Order

		if err := json.Unmarshal(data.([]byte), &order); err != nil {
			log.Printf("error during orders.shipped: %v", err)
		}
		log.Printf("order %d status = %v", order.Id, order.Status)

		// update order status
		err := os.repo.UpdateProductOrder(order)
		if err != nil {
			log.Printf("error updating product order %d: %v", order.Id, err)
		}
	})
	if err != nil {
		log.Printf("unable to get message from orders.shipped: %v", err)
	}
}

func (os *OrderService) SubscribeOrdersCompleted() {
	_, err := os.pubsub.Subscribe("orders.completed", func(data any) {
		var order model.Order

		if err := json.Unmarshal(data.([]byte), &order); err != nil {
			log.Printf("error during orders.completed: %v", err)
		}
		log.Printf("order %d status = %v", order.Id, order.Status)

		// update order status
		err := os.repo.UpdateProductOrder(order)
		if err != nil {
			log.Printf("error updating product order %d: %v", order.Id, err)
		}
	})
	if err != nil {
		log.Printf("unable to get message from orders.completed: %v", err)
	}
}
