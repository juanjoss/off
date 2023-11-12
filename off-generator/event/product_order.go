package event

import (
	"bytes"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

var (
	ordersEndpoint = "http://orders:8080/api/orders"
)

type ProductOrder struct{}

type ProductOrderRequest struct {
	SsdId    int    `json:"ssd_id"`
	Barcode  string `json:"barcode"`
	Quantity int    `json:"quantity"`
}

func (pr *ProductOrder) Handle() {
	var request ProductOrderRequest

	/*
		Publish a message through a NATS sunject to get and user's ssd id and a product's barcode.
		The request is published to "orders.random" and the reply is received on "orders.new".
	*/
	nc, err := nats.Connect("nats://nats:" + os.Getenv("NATS_PORT"))
	if err != nil {
		log.Printf("unable to connect to NATS: %v", err)
	}

	sub, err := nc.SubscribeSync("orders.new")
	if err != nil {
		log.Fatal(err)
	}
	nc.Flush()

	nc.PublishRequest("orders.random", "orders.new", nil)

	var count int
	max := 500 * time.Millisecond
	start := time.Now()
	for time.Since(start) < max {
		msg, err := sub.NextMsg(1 * time.Second)
		if err != nil {
			break
		}

		if err := json.Unmarshal(msg.Data, &request); err != nil {
			log.Printf("unable to unmarshal message %s: %v", string(msg.Data), err)
		}
		count++

		if count >= 2 {
			break
		}
	}
	// generate quantity
	request.Quantity = rand.Intn(4) + 1

	sub.Unsubscribe()
	nc.Close()

	// If no user registrations happend, then there can't be a product order
	if request.SsdId == 0 {
		log.Println("cancelling event", pr.Type(), "due to unexisting user-registration events.")
		return
	}

	// Send request to orders service
	jsonData, err := json.Marshal(request)
	if err != nil {
		log.Printf("unable to marshal request: %v", err)
		return
	}

	res, err := http.Post(
		ordersEndpoint,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		log.Printf("unable to POST product order: %v", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Println("event", pr.Type(), "failed")
		return
	}
}

func (pr *ProductOrder) Type() string {
	return "product-order"
}
