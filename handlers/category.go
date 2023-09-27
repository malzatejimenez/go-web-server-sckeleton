package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"platzi/go/rest-ws/models"
	"platzi/go/rest-ws/repository"
	"platzi/go/rest-ws/server"
	"strconv"

	"github.com/gorilla/mux"
)

// InsertCategoryRequest is a struct that contains the request body for the InsertCategory method
type InsertCategoryRequest struct {
	Name string `json:"name"`
}

// InsertCategoryResponse is a struct that contains the response body for the InsertCategory method
type InsertCategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// InsertCategoryHandler is a function that handles the InsertCategory method
func InsertCategoryHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// create a new variable to store the request body
		var req InsertCategoryRequest

		// decode the request body into the req variable
		err := decode(r.Body, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Errorf("error decoding request: %v", err))
			return
		}

		// create a new category
		category := &models.Category{
			Name: req.Name,
		}

		// insert the category into the database
		id, err := repository.InsertCategory(r.Context(), category)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error inserting category"))
			return
		}

		// create a new response
		res := &InsertCategoryResponse{
			ID:   id,
			Name: category.Name,
		}

		// write the header
		w.WriteHeader(http.StatusCreated)

		// set the content type
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		err = json.NewEncoder(w).Encode(res)
	}
}

// GetCategoryByIdRequest is a struct that contains the request body for the GetCategory method
type GetCategoryByIdRequest struct {
	ID int64 `json:"id"`
}

// GetCategoryResponse is a struct that contains the response body for the GetCategory method
type GetCategoryResponse struct {
	Category *models.Category `json:"category"`
}

// GetCategoryByIdHandler is a function that handles the GetCategoryById method
func GetCategoryByIdHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get params from the url
		params := mux.Vars(r)

		// parse id to int64
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.New("You must provide a valid id"))
			return
		}

		// get the category from the database
		category, err := repository.GetCategoryById(r.Context(), id)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error getting category"))
			return
		}

		// check if the category is nil
		if category == nil {
			respondError(w, http.StatusNotFound, errors.New("category not found"))
			return
		}

		// define the response
		res := &GetCategoryResponse{
			Category: category,
		}

		// write the header
		w.WriteHeader(http.StatusOK)

		// set the content type
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		err = json.NewEncoder(w).Encode(res)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error encoding response"))
			return
		}
	}
}

// UpdateCategoryRequest is a struct that contains the request body for the UpdateCategory method
type UpdateCategoryRequest struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UpdateCategoryResponse is a struct that contains the response body for the UpdateCategory method
type UpdateCategoryResponse struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

// UpdateCategoryHandler is a function that handles the UpdateCategory method
func UpdateCategoryHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get params from the url
		params := mux.Vars(r)

		// parse id to int64
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.New("You must provide a valid id"))
			return
		}

		// create a new variable to store the request body
		var req UpdateCategoryRequest

		// decode the request body into the req variable
		err = decode(r.Body, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, fmt.Errorf("error decoding request: %v", err))
			return
		}

		// get the category from the database
		category, err := repository.GetCategoryById(r.Context(), id)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error getting category"))
			return
		}

		// check if the category is nil
		if category == nil {
			respondError(w, http.StatusNotFound, errors.New("category not found"))
			return
		}

		// check if the new name is different from the old one
		if category.Name == req.Name {
			respondError(w, http.StatusBadRequest, errors.New("the new name must be different from the old one"))
			return
		}

		// due the name is is a unique field in the database, we need to check if the new name is already in use and return a bad request status if it is
		cat, _ := repository.GetCategoryByName(r.Context(), req.Name)
		if cat.Name == req.Name {
			respondError(w, http.StatusBadRequest, errors.New("the new name is already in use"))
			return
		}

		// update the category
		category.Name = req.Name

		// update the category into the database
		err = repository.UpdateCategory(r.Context(), category)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error updating category"))
			return
		}

		// create a new response
		res := &UpdateCategoryResponse{
			ID:   category.Id,
			Name: category.Name,
		}

		// write the header
		w.WriteHeader(http.StatusOK)

		// set the content type
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		err = json.NewEncoder(w).Encode(res)
	}
}

// DeleteCategoryResponse is a struct that contains the response body for the DeleteCategory method
type DeleteCategoryResponse struct {
	ID int64 `json:"id"`
}

// DeleteCategoryHandler is a function that handles the DeleteCategory method
func DeleteCategoryHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get params from the url
		params := mux.Vars(r)

		// parse id to int64
		id, err := strconv.ParseInt(params["id"], 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.New("The id must be a valid number"))
			return
		}

		// get the category from the database
		category, err := repository.GetCategoryById(r.Context(), id)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error getting category"))
			return
		}

		// check if the category is nil
		if category.Id == 0 {
			respondError(w, http.StatusNotFound, errors.New("category not found"))
			return
		}

		// delete the category from the database
		err = repository.DeleteCategory(r.Context(), id)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error deleting category"))
			return
		}

		// create a new response
		res := &DeleteCategoryResponse{
			ID: category.Id,
		}

		// write the header
		w.WriteHeader(http.StatusOK)

		// set the content type
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		json.NewEncoder(w).Encode(res)
	}
}

// ListCategoriesResponse is a struct that contains the response body for the ListCategories method
type ListCategoriesResponse struct {
	Categories []*models.Category `json:"categories"`
	Total      int64              `json:"total"`
}

// ListCategoriesHandler is a function that handles the ListCategories method
func ListCategoriesHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get query params
		page, err := strconv.ParseInt(r.URL.Query().Get("page"), 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.New("The page must be a valid number"))
			return
		}

		rowsPerPage, err := strconv.ParseInt(r.URL.Query().Get("rowsPerPage"), 10, 64)
		if err != nil {
			respondError(w, http.StatusBadRequest, errors.New("The rowsPerPage must be a valid number"))
			return
		}

		// list categories from the database
		categories, total, err := repository.ListCategories(r.Context(), page, rowsPerPage)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("error listing categories"))
			return
		}

		// create a new response
		res := &ListCategoriesResponse{
			Categories: categories,
			Total:      total,
		}

		// write the header
		w.WriteHeader(http.StatusOK)

		// set the content type
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		json.NewEncoder(w).Encode(res)
	}
}
