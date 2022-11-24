package main

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

var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

type Flow struct {
	Message   []string `json:"message"`
	NextState string   `json:"next_state"`
}

type Flows map[string]Flow

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

func (bot *Bot) AddIntent(title string, message []string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	bot.Intents[title] = message
	bot.UpdateBot()
	infoLog.Println("Intent added successfully")
}

func (bot *Bot) DeleteIntent(title string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bots exist")
	}
	delete(bot.Intents, title)
	bot.UpdateBot()
	infoLog.Println("Intent has been deleted")
}

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

func (bot *Bot) AddFlow(node string, flow Flow) {
	if bot.Id == 0 {
		errorLog.Fatal("No bot exists")
	}
	bot.Flows[node] = flow
	bot.UpdateBot()
	infoLog.Println("Flow added successfully")
}

func (bot *Bot) DeleteFlow(title string) {
	if bot.Id == 0 {
		errorLog.Fatal("No bots exist")
	}
	delete(bot.Flows, title)
	bot.UpdateBot()
	infoLog.Println("Flow has been deleted")
}

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

func fileChecker(fileName string) bool {
	_, error := os.Stat(fileName)

	// check if error is "file not exists"
	if !os.IsNotExist(error) {
		return true
	} else {
		return false
	}
}
