package repository

import (
	"context"
	"platzi/go/rest-ws/models"
)

// Repository interface is an interface that defines the methods that the repository should implement
type Repository interface {
	Close() error
	InsertUser(ctx context.Context, user *models.User) error
	GetUserById(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	InsertCategory(ctx context.Context, category *models.Category) (int64, error)
	GetCategoryById(ctx context.Context, id int64) (*models.Category, error)
	GetCategoryByName(ctx context.Context, name string) (*models.Category, error)
	UpdateCategory(ctx context.Context, category *models.Category) error
	DeleteCategory(ctx context.Context, id int64) error
	ListCategories(ctx context.Context, page, rowsPerPage int64) ([]*models.Category, int64, error)
}

// define a variable to store the implementation
var implementation Repository

// SetUserRepository is a setter for the implementation
func SetRepository(repo Repository) {
	implementation = repo
}

// Close is a function that calls the Close method of the implementation
func Close() error {
	return implementation.Close()
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

// InsertCategory is a function that calls the InsertCategory method of the implementation
func InsertCategory(ctx context.Context, category *models.Category) (int64, error) {
	return implementation.InsertCategory(ctx, category)
}

// GetCategoryById is a function that calls the GetCategoryById method of the implementation
func GetCategoryById(ctx context.Context, id int64) (*models.Category, error) {
	return implementation.GetCategoryById(ctx, id)
}

// GetCategoryByName is a function that calls the GetCategoryByName method of the implementation
func GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {
	return implementation.GetCategoryByName(ctx, name)
}

// UpdateCategory is a function that calls the UpdateCategory method of the implementation
func UpdateCategory(ctx context.Context, category *models.Category) error {
	return implementation.UpdateCategory(ctx, category)
}

// DeleteCategory is a function that calls the DeleteCategory method of the implementation
func DeleteCategory(ctx context.Context, id int64) error {
	return implementation.DeleteCategory(ctx, id)
}

// ListCategories is a function that calls the ListCategories method of the implementation
func ListCategories(ctx context.Context, page, rowsPerPage int64) ([]*models.Category, int64, error) {
	return implementation.ListCategories(ctx, page, rowsPerPage)
}
