package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"platzi/go/rest-ws/models"
	"platzi/go/rest-ws/repository"
	"platzi/go/rest-ws/server"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
)

// SignUpLoginRequest is a struct that represents the request of the SignUpHandler
type SignUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SignUpResponse is a struct that represents the response of the SignUpHandler
type SignUpResponse struct {
	Id    string `json:"id"`
	Email string `json:"email"`
}

// LoginResponse is a struct that represents the response of the LoginHandler
type LoginResponse struct {
	Token string `json:"token"`
}

// SignUpHandler is a function that handles the sign up of a user
func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// define a variable to decode the request into
		var req SignUpLoginRequest

		// decode the request
		err := decode(r.Body, &req)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		// generate the id
		id, err := ksuid.NewRandom()
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// create the user
		user := &models.User{
			Id:       id.String(),
			Email:    req.Email,
			Password: string(hashedPassword),
		}

		// insert the user
		err = repository.InsertUser(r.Context(), user)
		if err != nil {
			respondError(w, http.StatusInternalServerError, err)
			return
		}

		// create the response
		resp := SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		}

		// set the header
		w.Header().Set("Content-Type", "application/json")

		// set the status code
		w.WriteHeader(http.StatusCreated)

		// encode the response
		json.NewEncoder(w).Encode(resp)

	}
}

// LoginHandler is a function that handles the login of a user
func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// define a variable to decode the request into
		var request = SignUpLoginRequest{}

		// decode the request
		err := decode(r.Body, &request)
		if err != nil {
			respondError(w, http.StatusBadRequest, err)
			return
		}

		// get the user from the database
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		// check if the user is nil
		if user == nil {
			respondError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
			return
		}

		// compare the password
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password))
		if err != nil {
			respondError(w, http.StatusUnauthorized, errors.New("invalid credentials"))
			return
		}

		// create the claims
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
			},
		}

		// generate the token
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		// sign the token
		signedToken, err := token.SignedString([]byte(s.Config().JwtSecret))
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		// create the response
		resp := LoginResponse{
			Token: signedToken,
		}

		// set the header
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		json.NewEncoder(w).Encode(resp)

	}
}

// MeHandler is a function that handles the me endpoint
func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get the token from the header
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))

		// parse the token
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JwtSecret), nil
		})

		// check if there is an error
		if err != nil {
			respondError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		// get the claims
		claims, ok := token.Claims.(*models.AppClaims)
		if !ok || !token.Valid {
			respondError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		// get the user from the database
		user, err := repository.GetUserById(r.Context(), claims.UserId)
		if err != nil {
			log.Println(err)
			respondError(w, http.StatusInternalServerError, errors.New("internal server error"))
			return
		}

		// check if the user is nil
		if user == nil {
			respondError(w, http.StatusUnauthorized, errors.New("invalid token"))
			return
		}

		// create the response
		resp := SignUpResponse{
			Id:    user.Id,
			Email: user.Email,
		}

		// set the header
		w.Header().Set("Content-Type", "application/json")

		// encode the response
		json.NewEncoder(w).Encode(resp)
	}
}
