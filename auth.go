// Golang SDK for Sarufi Conversational AI Platform
package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Application struct {
	BaseURL string
	Token   string
}

// GetToken() method to generate a token for the user.
// The received token will be saved to Application.Token
// otherwise it will return an error.
func (app *Application) GetToken(username, password string) error {
	app.BaseURL = BaseURL
	url := app.BaseURL + "users/login"

	params := map[string]string{
		"username": username,
		"password": password,
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	switch resp.StatusCode {
	case 200:
		var result struct {
			Message string `json:"message"`
			Token   string `json:"token"`
		}

		err := json.Unmarshal([]byte(body), &result)

		if err != nil {
			return err
		}

		app.Token = result.Token

		return nil

	case 404:
		var notFound NotFoundError
		err := json.Unmarshal([]byte(body), &notFound)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s", notFound.Message)

	case 401:
		var unauthorized Unauthorized
		err := json.Unmarshal([]byte(body), &unauthorized)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s", unauthorized.Message)

	default:
		return fmt.Errorf(string(body))
	}

}
