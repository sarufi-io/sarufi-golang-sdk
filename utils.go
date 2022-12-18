package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/pretty"
)

// For Error and Information Loggings
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

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
	Name               string              `json:"name"`
	Industry           string              `json:"industry"`
	Description        string              `json:"description"`
	VisibleOnCommunity bool                `json:"visible_on_community"`
	Intents            map[string][]string `json:"intents"`
	Flows              `json:"flows"`
	baseURL            string
	chatID             string
	token              string
}

// CreateIntents method to create a new intent.
// It accepts a string that contains the intents
// arranged in JSON format. It also acceps a
// json file (eg: intents.json)
func (bot *Bot) CreateIntents(intents string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	if fileChecker(intents) {
		infoLog.Println("Creating intents from file...")
		newIntents, err := ioutil.ReadFile(intents)
		if err != nil {
			errorLog.Fatal(err)
		}
		for k := range bot.Intents {
			delete(bot.Intents, k)
		}

		err = json.Unmarshal(newIntents, &bot.Intents)
		if err != nil {
			errorLog.Fatalf("Could not create new intents from file: %v", err)
		}

	} else {
		for k := range bot.Intents {
			delete(bot.Intents, k)
		}
		newIntents := []byte(intents)
		err := json.Unmarshal(newIntents, &bot.Intents)
		if err != nil {
			errorLog.Fatalf("Could not create new intents: %v", err)
		}
	}
	bot.UpdateBot()
	infoLog.Println("Intents created successfully")
}

// AddIntent method to add a new intent. It accepts a string
// that will be the title of the intent, and a slice of message
// strings that will be the intent content.
func (bot *Bot) AddIntent(title string, message []string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	bot.Intents[title] = message
	bot.UpdateBot()
	infoLog.Println("Intent added successfully")
}

// DeleteIntent method to delete a specific intent.
// The accepted string is the title of the intent.
// If the title does not exist no error will be displayed.
func (bot *Bot) DeleteIntent(title string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bots exist")
	}
	delete(bot.Intents, title)
	bot.UpdateBot()
	infoLog.Println("Intent has been deleted")
}

// CreateFlows method to create a new flow.
// It accepts a string that contains the flows
// arranged in JSON format. It also acceps a
// json file (eg: flows.json)
func (bot *Bot) CreateFlows(flows string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	if fileChecker(flows) {
		infoLog.Println("Creating flows from file...")
		newFlows, err := ioutil.ReadFile(flows)
		if err != nil {
			errorLog.Fatal(err)
		}
		for k := range bot.Flows {
			delete(bot.Flows, k)
		}

		err = json.Unmarshal(newFlows, &bot.Flows)
		if err != nil {
			errorLog.Fatalf("Could not create new flows from file: %v", err)
		}
	} else {
		for k := range bot.Flows {
			delete(bot.Flows, k)
		}
		newFlows := []byte(flows)

		err := json.Unmarshal(newFlows, &bot.Flows)
		if err != nil {
			errorLog.Fatalf("Could not create new flows: %v", err)
		}
	}
	bot.UpdateBot()
	infoLog.Println("Flows created successfully")
}

// AddFlow method to add a new flow. It accepts a string
// that will be the title of the flow, and a type flow
// that will be the flow content.
func (bot *Bot) AddFlow(node string, flow Flow) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	bot.Flows[node] = flow
	bot.UpdateBot()
	infoLog.Println("Flow added successfully")
}

// DeleteFlow method to delete a specific flow.
// The accepted string is the title of the flow.
// If the title does not exist no error will be displayed.
func (bot *Bot) DeleteFlow(title string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bots exist")
	}
	delete(bot.Flows, title)
	bot.UpdateBot()
	infoLog.Println("Flow has been deleted")
}

// To get bot responses from the API. It accepts message
// which should be a string. And a messageType which should
// also be a string.
func (bot *Bot) Respond(message, messageType string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	infoLog.Print("Getting Response...")
	url := bot.baseURL + "conversation/"

	params := map[string]interface{}{
		"chat_id":      bot.chatID,
		"bot_id":       bot.Id,
		"message":      message,
		"message_type": messageType,
		"channel":      "general",
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		errorLog.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		errorLog.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	bearer := fmt.Sprintf("Bearer %s", bot.token)
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println(string(pretty.Pretty(body)))

}

// To get the current and next state of the chat. It
// accepts no parameters. It will display the JSON
// response from the API.
func (bot *Bot) ChatState() {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	infoLog.Print("Checking state...")
	url := bot.baseURL + "conversation/allchannels/status"
	params := map[string]interface{}{
		"chat_id": bot.chatID,
		"bot_id":  bot.Id,
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		errorLog.Fatal(err)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonParams))

	if err != nil {
		errorLog.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")
	bearer := fmt.Sprintf("Bearer %s", bot.token)
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println(string(pretty.Pretty(body)))
}

// A helper function to check if a file exists. It
// accepts a filename string and returns a bool.
func fileChecker(fileName string) bool {

	if len(fileName) > 50 {
		return false
	}
	_, error := os.Stat(fileName)

	// check if error is "file not exists"
	if !os.IsNotExist(error) {
		return true
	} else {
		return false
	}
}
