package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"platzi/go/rest-ws/models"
	"time"
)

// InsertCategory is a method that inserts a new category into the database
func (r *PostgresRepository) InsertCategory(ctx context.Context, category *models.Category) (int64, error) {
	// create the query
	query := `INSERT INTO categories(name) VALUES($1) RETURNING id`

	// create a new row
	row := r.db.QueryRowContext(ctx, query, category.Name)

	// create a new variable to store the id
	var id int64

	// scan the id
	err := row.Scan(&id)
	if err != nil {
		return 0, err
	}

	// return the id
	return id, nil
}

// GetCategoryById is a method that returns a category by its id
func (r *PostgresRepository) GetCategoryById(ctx context.Context, id int64) (*models.Category, error) {

	// define the query
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE id = $1`

	// execute the query
	rows, err := r.db.QueryContext(ctx, query, id)

	// check if there was an error
	if err != nil {
		return nil, fmt.Errorf("error getting category at GetUserByName: %v", err)
	}

	// return the category
	return extractCategoryFromResult(rows)
}

// GetCategoryByName is a method that returns a category by its id
func (r *PostgresRepository) GetCategoryByName(ctx context.Context, name string) (*models.Category, error) {

	// define the query
	query := `SELECT id, name, created_at, updated_at FROM categories WHERE name = $1`

	// execute the query
	rows, err := r.db.QueryContext(ctx, query, name)

	// check if there was an error
	if err != nil {
		return nil, fmt.Errorf("error getting category at GetUserByName: %v", err)
	}

	// return the category
	return extractCategoryFromResult(rows)
}

// extractCategoryFromResult is a function that extracts a category from a result
func extractCategoryFromResult(rows *sql.Rows) (*models.Category, error) {
	// define a defer to close the rows
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Printf("error closing rows at GetCategoryById: %v", err)
		}
	}()

	// define the category
	var category = models.Category{}

	// iterate over the rows
	for rows.Next() {
		// scan the row into the category
		err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)

		// check if there was an error scanning the row
		if err != nil {
			return nil, fmt.Errorf("error scanning category row at GetCategoryById: %v", err)
		}
	}

	// check if there was an error iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating over rows: %v", err)
	}

	// return the category
	return &category, nil
}

// UpdateCategory is a method that updates a category
func (r *PostgresRepository) UpdateCategory(ctx context.Context, category *models.Category) error {
	// define the query
	query := `UPDATE categories SET name = $1, updated_at = $2 WHERE id = $3`

	// define the updated at
	updatedAt := time.Now()

	// execute the query
	_, err := r.db.ExecContext(ctx, query, category.Name, updatedAt, category.Id)

	// check if there was an error
	if err != nil {
		return fmt.Errorf("error updating category at UpdateCategory: %v", err)
	}

	// return nil
	return nil
}

// DeleteCategory is a method that deletes a category
func (r *PostgresRepository) DeleteCategory(ctx context.Context, id int64) error {
	// define the query
	query := `DELETE FROM categories WHERE id = $1`

	// execute the query
	result, err := r.db.ExecContext(ctx, query, id)

	// check if there was an error
	if err != nil {
		return fmt.Errorf("error deleting category at DeleteCategory: %v", err)
	}

	// validate the result to see if the category was deleted
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error getting rows affected at DeleteCategory: %v", err)
	}

	// check if the category was deleted
	if rowsAffected == 0 {
		return fmt.Errorf("error deleting category at DeleteCategory: %v", err)
	}

	// return nil
	return nil
}

// ListCategories is a method that returns a list of categories
// Returns a list of categories and the total number of categories
func (r *PostgresRepository) ListCategories(ctx context.Context, page, rowsPerPage int64) ([]*models.Category, int64, error) {
	// define the query
	query := `SELECT id, name, created_at, updated_at FROM categories ORDER BY id LIMIT $1 OFFSET $2`

	// execute the query
	rows, err := r.db.QueryContext(ctx, query, rowsPerPage, (page-1)*rowsPerPage)

	// check if there was an error
	if err != nil {
		return nil, 0, fmt.Errorf("error getting categories at ListCategories: %v", err)
	}

	// define a defer to close the rows
	defer func() {
		err := rows.Close()
		if err != nil {
			fmt.Printf("error closing rows at ListCategories: %v", err)
		}
	}()

	// define the categories
	categories := make([]*models.Category, 0)

	// iterate over the rows
	for rows.Next() {
		// define the category
		var category = models.Category{}

		// scan the row into the category
		err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)

		// check if there was an error scanning the row
		if err != nil {
			return nil, 0, fmt.Errorf("error scanning category row at ListCategories: %v", err)
		}

		// append the category to the list of categories
		categories = append(categories, &category)
	}

	// check if there was an error iterating over the rows
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating over rows: %v", err)
	}

	// define the query to get the total number of categories
	query = `SELECT COUNT(*) FROM categories`

	// execute the query
	row := r.db.QueryRowContext(ctx, query)

	// define the total number of categories
	var total int64

	// scan the row into the total
	err = row.Scan(&total)

	// check if there was an error scanning the row
	if err != nil {
		return nil, 0, fmt.Errorf("error scanning total row at ListCategories: %v", err)
	}

	// return the categories and the total number of categories
	return categories, total, nil
}
