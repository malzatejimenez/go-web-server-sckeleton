package main

import (
	"context"
	"log"
	"os"
	"platzi/go/rest-ws/handlers"
	"platzi/go/rest-ws/middlewares"
	"platzi/go/rest-ws/server"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	// Get environment variables
	PORT := os.Getenv("PORT")
	JWT_SECRET := os.Getenv("JWT_SECRET")
	DATABASE_URL := os.Getenv("DATABASE_URL")

	// Create new server config
	config := &server.Config{
		Port:        PORT,
		JwtSecret:   JWT_SECRET,
		DatabaseURL: DATABASE_URL,
	}

	// Create new server
	s, err := server.NewServer(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	// Start server
	err = s.Start(BindRoutes)
	if err != nil {
		log.Fatal(err)
	}

}

// BindRoutes binds all routes to the router
func BindRoutes(s server.Server, r *mux.Router) {
	// Use CheckAuthMiddleware to check if the user is authenticated
	r.Use(middlewares.CheckAuthMiddleware(s))

	// Bind home handler
	r.HandleFunc("/", handlers.HomeHandler(s)).Methods("GET")

	// Bind Signup handler
	r.HandleFunc("/signup", handlers.SignUpHandler(s)).Methods("POST")

	// Bind Login handler
	r.HandleFunc("/login", handlers.LoginHandler(s)).Methods("POST")

	// Bind Me handler
	r.HandleFunc("/me", handlers.MeHandler(s)).Methods("GET")

	// Bind ListCategories handler
	r.HandleFunc("/categories", handlers.ListCategoriesHandler(s)).Methods("GET")

	// Bind InsertCategory handler
	r.HandleFunc("/categories", handlers.InsertCategoryHandler(s)).Methods("POST")

	// Bind GetCategoryById handler
	r.HandleFunc("/categories/{id:[0-9]+}", handlers.GetCategoryByIdHandler(s)).Methods("GET")

	// Bind UpdateCategory handler
	r.HandleFunc("/categories/{id:[0-9]+}", handlers.UpdateCategoryHandler(s)).Methods("PUT")

	// Bind DeleteCategory handler
	r.HandleFunc("/categories/{id:[0-9]+}", handlers.DeleteCategoryHandler(s)).Methods("DELETE")

}
