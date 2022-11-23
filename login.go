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
	if err != nil {
		return nil, fmt.Errorf("login: request execute: %w", err)
	}
	var loginResponse LoginResponse
	err = parseResponse(resp, &loginResponse)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	return &loginResponse, nil
}
