package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/juanjoss/off-users-service/ports"
)

const apiPrefix = "/api"

type Server struct {
	service ports.UserService
	router  *mux.Router
	port    string
}

func NewServer(us ports.UserService) *Server {
	return &Server{
		service: us,
		router:  mux.NewRouter().PathPrefix(apiPrefix).Subrouter(),
		port:    os.Getenv("APP_PORT"),
	}
}

func (s *Server) ListenAndServe() {
	// register http routes
	s.router.HandleFunc("/register", s.register).Methods(http.MethodPost)
	s.router.HandleFunc("/health", s.healthcheck).Methods(http.MethodGet)

	// register NATS subscriptions
	s.service.SubscribeOrdersNew()
	s.service.SubscribeOrdersCompleted()

	// create logging and recovery middlewares
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.router)
	recoveryRouter := handlers.RecoveryHandler()(loggedRouter)

	// run the http server
	log.Printf("users service running on port %s", s.port)
	log.Fatal(
		http.ListenAndServe(
			fmt.Sprintf(":%s", s.port),
			recoveryRouter,
		),
	)
}

func (s *Server) register(w http.ResponseWriter, r *http.Request) {
	var request ports.RegisterRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = s.service.Register(request)
	if err != nil {
		log.Println(err)
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
