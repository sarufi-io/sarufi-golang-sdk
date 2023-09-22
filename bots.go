package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
)

// CreateIntents method to create a new intent.
// It accepts a string that contains the intents
// arranged in JSON format. It also acceps a
// json file (eg: intents.json).
// It will return an error if any.
func (bot *Bot) CreateIntents(intents string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}

	if fileChecker(intents) {
		newIntents, err := ioutil.ReadFile(intents)
		if err != nil {
			return err
		}
		for k := range bot.Intents {
			delete(bot.Intents, k)
		}

		err = json.Unmarshal(newIntents, &bot.Intents)
		if err != nil {
			return err
		}

	} else {
		for k := range bot.Intents {
			delete(bot.Intents, k)
		}
		newIntents := []byte(intents)
		err := json.Unmarshal(newIntents, &bot.Intents)
		if err != nil {
			return err
		}
	}

	return nil
}

// AddIntent method to add a new intent. It accepts a string
// that will be the title of the intent, and a slice of message
// strings that will be the intent content.
func (bot *Bot) AddIntent(title string, message []string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	bot.Intents[title] = message
	return nil
}

// DeleteIntent method to delete a specific intent.
// The accepted string is the title of the intent.
// If the title does not exist no error will be displayed.
func (bot *Bot) DeleteIntent(title string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	delete(bot.Intents, title)
	return nil
}

// CreateFlows method to create a new flow.
// It accepts a string that contains the flows
// arranged in JSON format. It also acceps a
// json file (eg: flows.json)
func (bot *Bot) CreateFlows(flows string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	if fileChecker(flows) {
		newFlows, err := ioutil.ReadFile(flows)
		if err != nil {
			return err
		}
		for k := range bot.Flows {
			delete(bot.Flows, k)
		}

		err = json.Unmarshal(newFlows, &bot.Flows)
		if err != nil {
			return err
		}
	} else {
		for k := range bot.Flows {
			delete(bot.Flows, k)
		}
		newFlows := []byte(flows)

		err := json.Unmarshal(newFlows, &bot.Flows)
		if err != nil {
			return err
		}
	}
	return nil
}

// AddFlow method to add a new flow. It accepts a string
// that will be the title of the flow and an interface.
// To allow ability to add choices, flow has been made an
// interface.
// See: https://docs.sarufi.io/docs/Getting%20started%20/chatbots-addons#handling-choices
func (bot *Bot) AddFlow(node string, flow interface{}) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	bot.Flows[node] = flow
	return nil
}

// DeleteFlow method to delete a specific flow.
// The accepted string is the title of the flow.
// If the title does not exist no error will be displayed.
func (bot *Bot) DeleteFlow(title string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	delete(bot.Flows, title)
	return nil
}

// To get bot responses from the API. It accepts a message to
// responded to and a channel. The default channel is 'general'
// See: https://neurotech-africa.stoplight.io/docs/sarufi/4a3ab3e807c34-handle-conversation
func (bot *Bot) Respond(message, channel string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	url := baseURL + "conversation/"

	if bot.ChatID == "" {
		bot.ChatID = uuid.New().String()
	}

	params := map[string]interface{}{
		"chat_id":      bot.ChatID,
		"bot_id":       bot.Id,
		"message":      message,
		"message_type": "text",
		"channel":      channel,
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		return err
	}
	statusCode, body, err := makeRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		return err
	}
	switch statusCode {
	case 200:
		if bot.ModelName == "" {
			if err := json.Unmarshal(body, &bot.ConversationWithKnowledge); err != nil {
				return err
			}
			return nil
		}

		if err := json.Unmarshal(body, &bot.Conversation); err != nil {
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

// To get the current and next state of the chat. It
// accepts no parameters. It will display the JSON
// response from the API.
func (bot *Bot) ChatState() error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	url := baseURL + "conversation/status"
	params := map[string]interface{}{
		"chat_id": bot.ChatID,
		"bot_id":  bot.Id,
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		return nil
	}
	statusCode, body, err := makeRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		return nil
	}
	switch statusCode {
	case 200:
		if err := json.Unmarshal(body, &bot.Conversation); err != nil {
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

// A method to predict the intent of a particular message. It will return
// an error if any. The result will be stored at the bot.Prediction field.
func (bot *Bot) Predict(message string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	bot.Prediction.Message = message
	url := baseURL + "predict/intent"
	params := map[string]interface{}{
		"message": bot.Prediction.Message,
		"bot_id":  bot.Id,
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		return nil
	}
	statusCode, body, err := makeRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		return nil
	}
	switch statusCode {
	case 200:
		if err := json.Unmarshal(body, &bot.Prediction); err != nil {
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

// A method to get all users communicating with the bot.
// A list of chat users will be stored at bot.ChatUsers field.
func (bot *Bot) GetChatUsers() error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	url := fmt.Sprintf("%schatbot/%d/users", baseURL, bot.Id)
	statusCode, body, err := makeRequest("GET", url, nil)

	if err != nil {
		return nil
	}
	switch statusCode {
	case 200:
		if err := json.Unmarshal(body, &bot.ChatUsers); err != nil {
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

// Get conversation history of a particular ChatID. The result
// will be stored at the bot.ConversationHistory.
func (bot *Bot) GetChatHistory(chatID string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No bot exists")
	}
	url := fmt.Sprintf("%sconversation/history/%d/%s", baseURL, bot.Id, chatID)
	statusCode, body, err := makeRequest("GET", url, nil)

	if err != nil {
		return nil
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
