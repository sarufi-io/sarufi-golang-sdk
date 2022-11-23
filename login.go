package sarufi

import (
	"context"
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
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	req, err := createRequest(ctx, method, url, request, headers)
	if err != nil {
		return nil, fmt.Errorf("login: %w", err)
	}
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
