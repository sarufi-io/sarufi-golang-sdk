package main

import (
	"context"
	"fmt"
	"github.com/sarufi-io/sarufi-golang-sdk"
	"net/http"
	"os"
)

func main() {

	client := sarufi.NewClient(
		sarufi.WithHttpClient(http.DefaultClient),
		sarufi.WithLogger(os.Stdout),
		sarufi.WithDebug(true))

	// Register a new user
	registerRequest := &sarufi.RegisterRequest{
		Username: "",
		Password: "",
	}

	registerResponse, err := client.Register(context.Background(), registerRequest)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	fmt.Printf("login response: %+v\n", registerResponse)

	// Login
	loginRequest := &sarufi.LoginRequest{
		Username: "",
		Password: "",
	}

	loginResponse, err := client.Login(context.Background(), loginRequest)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	fmt.Printf("login response: %+v\n", loginResponse)

}
