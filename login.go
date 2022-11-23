package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	LoginResponse struct {
		Message string `json:"message,omitempty"`
		Token   string `json:"token,omitempty"`
	}
)

func login(ctx context.Context, client *http.Client, url, method string, request *LoginRequest) (*LoginResponse, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("login: marshal request: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("login: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("login: request execute: %w", err)
	}
	// if the code is 422, it means the request was invalid
	// and we can unmarshal the response into a validation error
	// if the code is 200 unmarshal the response into a login response
	// if the code is anything else,  stringfy the response body and return
	// it as an error
	if resp.StatusCode == http.StatusOK {
		var loginResponse LoginResponse
		if err := json.NewDecoder(resp.Body).Decode(&loginResponse); err != nil {
			return nil, fmt.Errorf("login: decode login response: %w", err)
		}
		return &loginResponse, nil
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		var validationError ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&validationError); err != nil {
			return nil, fmt.Errorf("login: decode validation error: %w", err)
		}
		return nil, fmt.Errorf("login: %w", &validationError)
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("login: read response body: %w", err)
		}
		return nil, fmt.Errorf("login: %s", responseBody.String())
	}
}
