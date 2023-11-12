package pubsub

import (
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type PubSub struct {
	Conn *nats.EncodedConn
}

func NewPubSub() *PubSub {
	nc, err := nats.Connect("nats://nats:" + os.Getenv("NATS_PORT"))
	if err != nil {
		log.Printf("unable to connect to NATS: %v", err)
	}

	c, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		log.Printf("unable to create NATS JSON encoded connection: %v", err)
	}

	pubsub := &PubSub{
		Conn: c,
	}

	return pubsub
}

func (ps *PubSub) Publish(subject string, data any) error {
	return ps.Conn.Publish(subject, data)
}

func (ps *PubSub) Subscribe(subject string, handler func(data any)) (*nats.Subscription, error) {
	s, err := ps.Conn.Subscribe(subject, func(msg *nats.Msg) {
		handler(msg.Data)
	})
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (ps *PubSub) CloseConn() {
	ps.Conn.Close()
}
