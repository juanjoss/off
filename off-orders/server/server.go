package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/juanjoss/off-orders-service/ports"
)

const apiPrefix = "/api"

type Server struct {
	service ports.OrderService
	router  *mux.Router
	port    string
}

func NewServer(service ports.OrderService) *Server {
	return &Server{
		service: service,
		router:  mux.NewRouter().PathPrefix(apiPrefix).Subrouter(),
		port:    os.Getenv("APP_PORT"),
	}
}

func (s *Server) ListenAndServe() {
	// register http routes
	s.router.HandleFunc("/orders", s.createProductOrder).Methods(http.MethodPost)
	s.router.HandleFunc("/health", s.healthcheck).Methods(http.MethodGet)

	// register NATS subscriptions
	s.service.SubscribeOrdersNew()
	s.service.SubscribeOrdersShipped()
	s.service.SubscribeOrdersCompleted()

	// create logging and recovery middlewares
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.router)
	recoveryRouter := handlers.RecoveryHandler()(loggedRouter)

	// run the http server
	log.Printf("orders service running on port %s", s.port)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", s.port),
			recoveryRouter,
		),
	)
}

func (s *Server) createProductOrder(w http.ResponseWriter, r *http.Request) {
	var request ports.CreateProductOrderRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Printf("unable to decode request: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.service.CreateProductOrder(request)
	if err != nil {
		log.Printf("unable to create product order: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(http.StatusOK)
}

func (s *Server) healthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
