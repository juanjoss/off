package ports

import (
	"github.com/juanjoss/off-orders-service/model"
	"github.com/nats-io/nats.go"
)

/*
	Interfaces
*/

type OrderRepository interface {
	// Returns a random product's barcode. An error if it fails.
	GetRandomProductBarcode() (string, error)

	// Inserts a model.Order into the repository. An error if it fails.
	CreateProductOrder(model.Order) (int, error)

	// Updates a product order status
	UpdateProductOrder(model.Order) error
}

type PubSub interface {
	/*
		Publish(subject string, message any).
		Publishes message to subject.
		An error on a failed publish.
	*/
	Publish(string, any) error

	/*
		Subscribe(subject string, handler func(any)).
		Subscribes to subject and registers its handler.
		Returns the subscription and an error if the subscription fails.
	*/
	Subscribe(string, func(any)) (*nats.Subscription, error)

	// Closes the connection to the message broker.
	CloseConn()
}

type OrderService interface {
	// Creates the order specified by CreateProductOrderRequest.
	CreateProductOrder(CreateProductOrderRequest) error

	// Subscribes to orders.random and publishes to order.new (generator response).
	SubscribeOrdersNew()

	// Subscribes to orders.random and updates order's status.
	SubscribeOrdersShipped()

	// Subscribes to orders.completed and updates order's status.
	SubscribeOrdersCompleted()
}

/*
	Service Models
*/

// Service model used to parse new order requests (api/orders).
type CreateProductOrderRequest struct {
	SsdId    int    `json:"ssd_id"`
	Barcode  string `json:"barcode"`
	Quantity int    `json:"quantity"`
}
