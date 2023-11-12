package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/juanjoss/off-notifications-service/model"
	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://nats:" + os.Getenv("NATS_PORT"))
	if err != nil {
		log.Printf("unable to connect to NATS: %v", err)
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Printf("unable to create NATS JSON encoded connection: %v", err)
	}

	/*
		Supplier workflow
	*/
	c.Subscribe("orders.pending", func(msg *nats.Msg) {
		var order model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("(supplier) error during order.pending: %v", err)
		}

		n := rand.Intn(4) + 1
		log.Printf("(supplier) order %v received and scheduled to ship in %d minutes", order.Id, n)
		time.Sleep(time.Duration(n) * time.Minute)

		order.Status = "shipped"

		if err := c.Publish("orders.shipped", order); err != nil {
			log.Printf("(supplier) unable to publish to orders.shipped: %v", err)
		}
		log.Printf("(supplier) order %v shipped", order.Id)
	})

	/*
		Delivery workflow
	*/
	c.Subscribe("orders.shipped", func(msg *nats.Msg) {
		var order model.Order
		if err := json.Unmarshal(msg.Data, &order); err != nil {
			log.Printf("(delivery) error during order.shipped: %v", err)
		}

		n := rand.Intn(4) + 1
		log.Printf("(delivery) delivering order %v and scheduled to arrive in %d minutes", order.Id, n)
		time.Sleep(time.Duration(n) * time.Minute)

		order.Status = "completed"

		if err := c.Publish("orders.completed", order); err != nil {
			log.Printf("(delivery) unable to publish to orders.completed: %v", err)
		}
		log.Printf("(delivery) order %v delivered", order.Id)
	})

	for {
	}
}
