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
	if !checkToken(app.Token) {
		return nil, fmt.Errorf("Error: No token available")
	}
	url := fmt.Sprintf("%schatbot/%d", app.BaseURL, id)

	statusCode, body, err := app.makeRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	switch statusCode {
	case 200:
		var bot *Bot
		if err := json.Unmarshal(body, &bot); err != nil {
			return nil, err
		}
		// chatID := uuid.New()
		// bot.chatID = chatID.String()
		return bot, nil
	case 404:
		var notFound NotFoundError
		if err := json.Unmarshal(body, &notFound); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", notFound.Error())
	case 401:
		var unauthorized Unauthorized
		if err := json.Unmarshal(body, &unauthorized); err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("Error %s", unauthorized.Error())
	default:
		return nil, fmt.Errorf(string(body))
	}
}

// GetBots() method returns a list of type Bot and an error
func (app *Application) GetBots() ([]Bot, error) {
	if !checkToken(app.Token) {
		return nil, fmt.Errorf("Error: No token available")
	}

	url := app.BaseURL + "chatbots"

	statusCode, body, err := app.makeRequest("GET", url, nil)
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
	if !checkToken(app.Token) {
		return nil, fmt.Errorf("Error: No token available")
	}

	url := app.BaseURL + "chatbot"

	params := map[string]interface{}{
		"name":                 name,
		"description":          description,
		"industry":             industry,
		"visible_on_community": visible,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		errorLog.Fatal(err)
	}

	statusCode, body, err := app.makeRequest("POST", url, bytes.NewBuffer(jsonParams))

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

	default:
		return nil, fmt.Errorf(string(body))
	}
}

// UpdateBot() method to update the bot. It accepts a parameter of
// type *Bot and will return  an error if any
func (app *Application) UpdateBot(bot *Bot) error {
	if !checkToken(app.Token) {
		return fmt.Errorf("Error: No token available")
	}

	jsonParams, err := json.Marshal(bot)

	if err != nil {
		return err
	}
	url := fmt.Sprintf("%schatbot/%d", app.BaseURL, bot.Id)
	statusCode, body, err := app.makeRequest("PUT", url, bytes.NewBuffer(jsonParams))
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

	default:
		return fmt.Errorf(string(body))
	}

}

// DeleteBot() will delete the bot of with provided ID.
// It will return an error if deletion was unsuccessful
func (app *Application) DeleteBot(id int) error {
	if !checkToken(app.Token) {
		return fmt.Errorf("Error: No token available")
	}

	url := fmt.Sprintf("%schatbot/%d", app.BaseURL, id)
	statusCode, body, err := app.makeRequest("DELETE", url, nil)
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

	default:
		return fmt.Errorf(string(body))
	}
}
