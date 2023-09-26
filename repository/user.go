package repository

import (
	"context"
	"platzi/go/rest-ws/models"
)

// UserRepository interface is an interface that defines the methods that the repository should implement
type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Close() error
}

// define a variable to store the implementation
var implementation UserRepository

// SetUserRepository is a setter for the implementation
func SetUserRepository(repo UserRepository) {
	implementation = repo
}

// InsertUser is a function that calls the InsertUser method of the implementation
func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

// GetUserById is a function that calls the GetUserById method of the implementation
func GetUserById(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserById(ctx, id)
}

// GetUserByEmail is a function that calls the GetUserByEmail method of the implementation
func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

// Close is a function that calls the Close method of the implementation
func Close() error {
	return implementation.Close()
}
