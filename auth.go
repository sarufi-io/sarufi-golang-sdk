// Golang SDK for Sarufi Conversational AI Platform
package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetToken() method to generate a token for the user.
// The received token will be saved to Application.Token
// otherwise it will return an error.
func (app *Application) GetToken(client_id, client_secret string) error {
	url := baseURL + "api/access_token"

	params := map[string]string{
		"client_id":     client_id,
		"client_secret": client_secret,
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
			Expires int    `json:"expires_in"`
			Token   string `json:"access_token"`
		}
		err := json.Unmarshal([]byte(body), &result)
		if err != nil {
			return err
		}
		token = result.Token
		return nil
	case 401:
		var unauthorized Unauthorized
		err := json.Unmarshal([]byte(body), &unauthorized)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		err := json.Unmarshal([]byte(body), &notFound)
		if err != nil {
			return err
		}
		return fmt.Errorf("Error %s", notFound.Error())
	case 422:
		var unprocessableEntity *UnprocessableEntity
		if err := json.Unmarshal(body, &unprocessableEntity); err != nil {
			return err
		}
		return fmt.Errorf("Error %s", unprocessableEntity.Error())
	case 500:
		return fmt.Errorf("Error status code 500: Internal Server Error")
	default:
		return fmt.Errorf(string(body))
	}
}
