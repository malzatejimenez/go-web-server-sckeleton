package postgres

import (
	"database/sql"
	"fmt"
	"platzi/go/rest-ws/repository"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository is a function that returns a new PostgresRepository
func NewPostgresRepository(url string) (repository.UserRepository, error) {
	// open the connection
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	// return the repository
	return &PostgresRepository{db}, nil
}

// Close is a method that closes the connection to the database
func (repo *PostgresRepository) Close() error {
	// close the connection
	err := repo.db.Close()

	// check if there was an error closing the connection
	if err != nil {
		return fmt.Errorf("error closing connection: %v", err)
	}

	// return nil as error
	return nil
}
