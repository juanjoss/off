package main

import (
	"github.com/juanjoss/off-orders-service/pubsub"
	"github.com/juanjoss/off-orders-service/repository"
	"github.com/juanjoss/off-orders-service/server"
	"github.com/juanjoss/off-orders-service/service"

	_ "github.com/juanjoss/off-orders-service/docs"
)

func main() {
	pr := repository.NewProductRepository()
	pps := pubsub.NewPubSub()
	ps := service.NewOrdersService(pr, pps)
	defer pps.CloseConn()
	server.NewServer(ps).ListenAndServe()
}
