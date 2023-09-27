package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"platzi/go/rest-ws/database/postgres"
	"platzi/go/rest-ws/repository"

	"github.com/gorilla/mux"
)

// Config is the server config struct
type Config struct {
	Port        string
	JwtSecret   string
	DatabaseURL string
}

// Server is the interface that all servers must implement
type Server interface {
	Config() *Config
}

// Broker is the server struct that implements the Server interface
type Broker struct {
	config *Config
	router *mux.Router
}

// Config returns the server config
func (b *Broker) Config() *Config {
	return b.config
}

// NewServer creates a new server instance
func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	// Validate config port is not empty
	if config.Port == "" {
		return nil, errors.New("port is required")
	}

	// Validate config JWTSecret is not empty
	if config.JwtSecret == "" {
		return nil, errors.New("jwt secret is required")
	}

	// Validate config DatabaseURL is not empty
	if config.DatabaseURL == "" {
		return nil, errors.New("database url is required")
	}

	// Create new broker
	broker := &Broker{
		config: config,
		router: mux.NewRouter(),
	}

	// Return broker and a nil error
	return broker, nil
}

// Start starts the server
func (b *Broker) Start(binder func(s Server, r *mux.Router)) error {
	// Inits the broker router
	b.router = mux.NewRouter()

	// Bind the router
	binder(b, b.router)

	// init repository
	repo, err := postgres.NewPostgresRepository(b.config.DatabaseURL)
	if err != nil {
		return fmt.Errorf("error initializing repository: %v", err)
	}

	// init abstract repository
	repository.SetRepository(repo)

	// Loging server start
	log.Println("Server started on port", b.config.Port)

	// Start the server
	if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
		log.Fatal("ListenAndServe", err)
		return err
	}

	// Return nil error
	return nil
}
