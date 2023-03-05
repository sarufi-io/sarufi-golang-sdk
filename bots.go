package sarufi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// Type Flow to create a new Flow XD.
// Message is a slice of strings.
// NextState is simply a string.
type Flow struct {
	Message   []string `json:"message"`
	NextState string   `json:"next_state"`
}

// Type Flows to map a string to the type Flow.
type Flows map[string]Flow

// Type Bot. All the fields are matched to the API
// JSON response.
type Bot struct {
	Id                 int                 `json:"id"`
	UserID             int                 `json:"user_id"`
	Username           string              `json:"user_name"`
	Name               string              `json:"name"`
	Industry           string              `json:"industry"`
	Description        string              `json:"description"`
	VisibleOnCommunity bool                `json:"visible_on_community"`
	Intents            map[string][]string `json:"intents"`
	Flows              Flows               `json:"flows"`
	ModelName          string              `json:"model_name"`
	CreatedAt          string              `json:"created_at"`
	UpdatedAt          string              `json:"updated_at"`
	chatID             string
}

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
		return fmt.Errorf("No boy exists")
	}
	bot.Intents[title] = message
	return nil
}

// DeleteIntent method to delete a specific intent.
// The accepted string is the title of the intent.
// If the title does not exist no error will be displayed.
func (bot *Bot) DeleteIntent(title string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No boy exists")
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
		return fmt.Errorf("No boy exists")
	}
	if fileChecker(flows) {
		infoLog.Println("Creating flows from file...")
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
// that will be the title of the flow, and a type flow
// that will be the flow content.
func (bot *Bot) AddFlow(node string, flow Flow) error {
	if bot.Id == 0 {
		return fmt.Errorf("No boy exists")
	}
	bot.Flows[node] = flow
	return nil
}

// DeleteFlow method to delete a specific flow.
// The accepted string is the title of the flow.
// If the title does not exist no error will be displayed.
func (bot *Bot) DeleteFlow(title string) error {
	if bot.Id == 0 {
		return fmt.Errorf("No boy exists")
	}
	delete(bot.Flows, title)
	return nil
}

// To get bot responses from the API. It accepts message
// which should be a string. And a messageType which should
// also be a string.
// func (bot *Bot) Respond(message, messageType string) error {
// 	if bot.Id == 0 {
// 		return fmt.Errorf("No boy exists")
// 	}
// 	url := BaseURL + "conversation/"
//
// 	params := map[string]interface{}{
// 		"chat_id":      bot.chatID,
// 		"bot_id":       bot.Id,
// 		"message":      message,
// 		"message_type": messageType,
// 		"channel":      "general",
// 	}
//
// 	jsonParams, err := json.Marshal(params)
//
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
//
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
//
// 	req.Header.Set("Content-Type", "application/json")
// 	bearer := fmt.Sprintf("Bearer %s", bot.token)
// 	req.Header.Set("Authorization", bearer)
//
// 	client := &http.Client{}
//
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
//
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	infoLog.Println(string(body))
// 	return nil
//
// }
//
// // To get the current and next state of the chat. It
// // accepts no parameters. It will display the JSON
// // response from the API.
// func (bot *Bot) ChatState() error {
// 	if bot.Id == 0 {
// 		errorLog.Fatal("No bot exists")
// 	}
// 	infoLog.Print("Checking state...")
// 	url := bot.baseURL + "conversation/allchannels/status"
// 	params := map[string]interface{}{
// 		"chat_id": bot.chatID,
// 		"bot_id":  bot.Id,
// 	}
//
// 	jsonParams, err := json.Marshal(params)
//
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))
//
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
//
// 	req.Header.Set("Content-Type", "application/json")
// 	bearer := fmt.Sprintf("Bearer %s", bot.token)
// 	req.Header.Set("Authorization", bearer)
//
// 	client := &http.Client{}
//
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	defer resp.Body.Close()
//
// 	body, err := io.ReadAll(resp.Body)
//
// 	if err != nil {
// 		errorLog.Fatal(err)
// 	}
// 	infoLog.Println(string(body))
// 	return nil
// }
