package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	// RegisterRequest is the request body for registering a new user
	RegisterRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// RegisterResponse is the response body for registering a new user
	RegisterResponse struct {
		Message string `json:"message,omitempty"`
	}
)

// register is a helper function that handles the request and response for registering a new user
// its a low level function that is used by higher level functions. RegisterResponse is returned
// unless there is an error or the response status code is http.StatusUnprocessableEntity
// where contents of ValidationError is returned. Else, the response body is returned as an error
func register(ctx context.Context, client *http.Client, url, method string,
	request *RegisterRequest) (*RegisterResponse, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("register: marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("register: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("register: request execute: %w", err)
	}
	var registerResponse RegisterResponse
	err = parseResponse(resp, &registerResponse)
	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return &registerResponse, nil
}
