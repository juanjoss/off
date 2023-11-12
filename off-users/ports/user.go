package ports

import (
	"github.com/juanjoss/off-users-service/model"
	"github.com/nats-io/nats.go"
)

/*
	Interfaces
*/

type UserRepository interface {
	// Returns a random SSD's id. An error if it fails.
	RandomSSDId() (int, error)

	// Inserts a user and its SSDs. An error if it fails.
	Register(RegisterRequest) error

	// Inserts a product into an user's SSD. An error if it fails.
	AddProductToSSD(int, string, int) error
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

type UserService interface {
	// Registers a new user and its SSDs specified by RegisterRequest.
	Register(RegisterRequest) error

	// Subscribes to orders.random and publishes to order.new (generator response).
	SubscribeOrdersNew()

	// Subscribes to orders.completed.
	SubscribeOrdersCompleted()
}

/*
	Service models
*/

// Service model used to parse user's registration requests (api/register).
type RegisterRequest struct {
	User *model.User  `json:"user"`
	SSDs []*model.SSD `json:"devices"`
}

// Service model used to parse a completed product's order (orders.completed subject).
type AddProductToSSD struct {
	SsdId    int    `json:"ssd_id"`
	Barcode  string `json:"product_id"`
	Quantity int    `json:"quantity"`
}
