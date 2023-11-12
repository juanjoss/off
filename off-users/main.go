package main

import (
	"github.com/juanjoss/off-users-service/pubsub"
	"github.com/juanjoss/off-users-service/repository"
	"github.com/juanjoss/off-users-service/server"
	"github.com/juanjoss/off-users-service/service"

	_ "github.com/juanjoss/off-users-service/docs"
)

func main() {
	ur := repository.NewUserRepository()
	ups := pubsub.NewPubSub()
	us := service.NewUserService(ur, ups)
	defer ups.CloseConn()
	server.NewServer(us).ListenAndServe()
}
