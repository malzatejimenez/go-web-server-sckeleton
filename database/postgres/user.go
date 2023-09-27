package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"platzi/go/rest-ws/models"
)

// InsertUser is a method that inserts a user into the database
func (r *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	// execute the query
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)

	// check if there was an error
	if err != nil {
		return fmt.Errorf("error inserting user: %v", err)
	}

	// return nil as error
	return nil
}

// GetUserById is a method that returns a user from the database
func (r *PostgresRepository) GetUserById(ctx context.Context, id string) (*models.User, error) {
	// execute the query
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE id = $1", id)

	// check if there was an error
	if err != nil {
		return nil, fmt.Errorf("error getting user at GetUserById: %v", err)
	}

	// get the user from the result
	user, err := extractUserFromResult(rows)
	if err != nil {
		return nil, fmt.Errorf("error getting user at GetUserById: %v", err)
	}

	// removing the password for security reasons
	user.Password = ""

	// return the user
	return user, nil
}

// GetUserByEmail is a method that returns a user from the database
func (r *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// execute the query
	rows, err := r.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)

	// check if there was an error
	if err != nil {
		return nil, fmt.Errorf("error getting user at GetUserByEmail: %v", err)
	}

	// return the user
	return extractUserFromResult(rows)
}

// extractUserFromResult is a function that extracts a user from a result
func extractUserFromResult(rows *sql.Rows) (*models.User, error) {
	// define a defer to close the rows
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Printf("error closing rows at GetUserByEmail: %v", err)
		}
	}()

	// define the user
	var user = models.User{}

	// iterate over the rows
	for rows.Next() {
		// scan the row into the user
		err := rows.Scan(&user.Id, &user.Email, &user.Password)

		// check if there was an error scanning the row
		if err != nil {
			return nil, fmt.Errorf("error scanning user row at GetUserByEmail: %v", err)
		}
	}

	// check if there was an error iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	// return the user
	return &user, nil
}
