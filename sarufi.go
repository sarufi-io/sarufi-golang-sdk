package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// GetBot() method to get a specific bot.
// It accepts the bot id (type int) as a parameter.
// Returns a pointer of type Bot and an error.
func (app *Application) GetBot(id int) (*Bot, error) {
	if !checkToken(token) {
		return nil, fmt.Errorf("Error: No token available")
	}
	url := fmt.Sprintf("%schatbot/%d", baseURL, id)

	statusCode, body, err := makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 200:
		var bot *Bot
		if err := json.Unmarshal(body, &bot); err != nil {
			return nil, err
		}
		return bot, nil
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", notFound.Error())
	case 422:
		var unprocessableEntity *UnprocessableEntity
		if err := json.Unmarshal(body, &unprocessableEntity); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unprocessableEntity.Error())
	case 500:
		return nil, fmt.Errorf("Error status code 500: Internal Server Error")
	default:
		return nil, fmt.Errorf(string(body))
	}
}

// GetBots() method returns a list of type Bot and an error
func (app *Application) GetAllBots() ([]Bot, error) {
	if !checkToken(token) {
		return nil, fmt.Errorf("Error: No token available")
	}

	url := baseURL + "chatbots"

	statusCode, body, err := makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	switch statusCode {
	case 200:
		var bots []Bot
		if err := json.Unmarshal(body, &bots); err != nil {
			return nil, err
		}
		return bots, nil
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", notFound.Error())
	case 422:
		var unprocessableEntity *UnprocessableEntity
		if err := json.Unmarshal(body, &unprocessableEntity); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unprocessableEntity.Error())
	case 500:
		return nil, fmt.Errorf("Error status code 500: Internal Server Error")
	default:
		return nil, fmt.Errorf(string(body))
	}

}

// CreateBot method to create a new bot. It accepts the following parameters:
// Name (type string) - which is the name of the bot
// Description (type string)- which is the description of the bot
// Industry (type string) - which is the related industry of the bot
// Visible (type bool) - allow it to be publicly available
// It returns a pointer to a new Bot and an error
func (app *Application) CreateBot(name, description, industry string, visible bool) (*Bot, error) {
	if !checkToken(token) {
		return nil, fmt.Errorf("Error: No token available")
	}

	url := baseURL + "chatbot"

	params := map[string]interface{}{
		"name":                 name,
		"description":          description,
		"industry":             industry,
		"visible_on_community": visible,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	statusCode, body, err := makeRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 200:
		var bot *Bot
		if err := json.Unmarshal(body, &bot); err != nil {
			return nil, err
		}
		return bot, nil
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", notFound.Error())
	case 422:
		var unprocessableEntity *UnprocessableEntity
		if err := json.Unmarshal(body, &unprocessableEntity); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unprocessableEntity.Error())
	case 500:
		return nil, fmt.Errorf("Error status code 500: Internal Server Error")
	default:
		return nil, fmt.Errorf(string(body))
	}
}

// UpdateBot() method to update the bot. It accepts a parameter of
// type *Bot and will return  an error if any
func (app *Application) UpdateBot(bot *Bot) error {
	if !checkToken(token) {
		return fmt.Errorf("Error: No token available")
	}

	jsonParams, err := json.Marshal(bot)

	if err != nil {
		return err
	}
	url := fmt.Sprintf("%schatbot/%d", baseURL, bot.Id)
	statusCode, body, err := makeRequest("PUT", url, bytes.NewBuffer(jsonParams))
	if err != nil {
		return err
	}

	switch statusCode {
	case 200:
		if err := json.Unmarshal(body, &bot); err != nil {
			return err
		}
		return nil
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return err
		}
		return fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
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

// DeleteBot() will delete the bot of with provided ID.
// It will return an error if deletion was unsuccessful
func (app *Application) DeleteBot(id int) error {
	if !checkToken(token) {
		return fmt.Errorf("Error: No token available")
	}

	url := fmt.Sprintf("%schatbot/%d", baseURL, id)
	statusCode, body, err := makeRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	switch statusCode {
	case 200:
		return nil
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return err
		}
		return fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
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

// Get user's information
func (app *Application) GetUser() (*User, error) {
	if !checkToken(token) {
		return nil, fmt.Errorf("Error: No token available")
	}

	url := fmt.Sprintf("%sapi/profile", baseURL)
	statusCode, body, err := makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 200:
		if err := json.Unmarshal(body, &app.User); err != nil {
			return nil, err
		}
		var user *User
		user = &app.User

		return user, nil
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unauthorized.Error())
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", notFound.Error())
	case 422:
		var unprocessableEntity *UnprocessableEntity
		if err := json.Unmarshal(body, &unprocessableEntity); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unprocessableEntity.Error())
	case 500:
		return nil, fmt.Errorf("Error status code 500: Internal Server Error")
	default:
		return nil, fmt.Errorf(string(body))
	}
}
