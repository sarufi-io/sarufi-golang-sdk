package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/tidwall/pretty"
)

var (
	infoLog  = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)
)

type Bot struct {
	Id                 int    `json:"id"`
	Name               string `json: "name"`
	Industry           string `json: "description"`
	Description        string `json: "description"`
	VisibleOnCommunity bool   `json: "visible_on_community"`
	Intents            struct {
		Intent map[string][]string
	}
	Flow struct {
		Message   map[string][]string
		NextState map[string]string
	}
	baseURL string
	token   string
}

func (bot *Bot) GetBot(id int) *Bot {
	infoLog.Printf("Getting bot with id: %d", id)
	url := fmt.Sprintf("%schatbot/%d", bot.baseURL, id)
	req, err := http.NewRequest("GET", url, nil)
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
	if err := json.Unmarshal(body, &bot); err != nil {
		errorLog.Fatal("Cannot Unmarshal JSON response.")
	}

	return bot
}

func (bot *Bot) GetBots() {
	infoLog.Print("Getting bots")
	url := bot.baseURL + "chatbots"
	req, err := http.NewRequest("GET", url, nil)
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

func (bot *Bot) CreateBot(name, description, industry string) *Bot {
	if bot.token == "" {
		errorLog.Fatal("Initialize the bot first!")
	}

	bot.Name = name
	bot.Description = description
	bot.Industry = industry
	infoLog.Print("Creating new bot...")
	url := bot.baseURL + "chatbot"
	params := map[string]interface{}{
		"name":                 bot.Name,
		"description":          bot.Description,
		"industry":             bot.Industry,
		"visible_on_community": bot.VisibleOnCommunity,
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
	if err := json.Unmarshal(body, &bot); err != nil {
		errorLog.Fatal("Cannot Unmarshal JSON response.")
	}

	return bot
}

func (bot *Bot) UpdateBot(id int) *Bot {
	infoLog.Print("Updating bot...")

	params := map[string]interface{}{
		"name":                 bot.Name,
		"description":          bot.Description,
		"industry":             bot.Industry,
		"visible_on_community": bot.VisibleOnCommunity,
		"intents":              bot.Intents,
		"flows":                bot.Flow,
	}

	jsonParams, err := json.Marshal(params)
	if err != nil {
		errorLog.Fatal(err)
	}
	url := fmt.Sprintf("%schatbot/%d", bot.baseURL, id)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonParams))
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
	return bot
}

func (bot *Bot) DeleteBot(id int) {
	infoLog.Printf("Deleting bot with id: %d", id)
	url := fmt.Sprintf("%schatbot/%d", bot.baseURL, id)
	req, err := http.NewRequest("DELETE", url, nil)
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
