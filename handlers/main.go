package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// decode is a function that decodes a request
func decode(body io.ReadCloser, v interface{}) error {
	// decode the request
	err := json.NewDecoder(body).Decode(&v)

	// check if there was an error
	if err != nil {
		return fmt.Errorf("error decoding request: %v", err)
	}

	// return nil as error
	return nil
}

// respondError is a function that responds with an error
func respondError(w http.ResponseWriter, code int, err error) {
	// write the header
	w.WriteHeader(code)

	// write the error
	w.Write([]byte(err.Error()))
}
