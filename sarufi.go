package sarufi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/tidwall/pretty"
)

// GetBot() method to get a specific bot.
// It accepts the bot id (type int) as a parameter.
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

	//	infoLog.Println(string(pretty.Pretty(body)))
	if err := json.Unmarshal(body, &bot); err != nil {
		errorLog.Fatal("Cannot Unmarshal JSON response.")
	}
	chatID := uuid.New()
	bot.chatID = chatID.String()
	infoLog.Printf("Fetched bot with id: %d", id)
	return bot

}

// GetBots() get a json response of all available bots and
// their details.
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

// CreateBot method to create a new bot. It accepts the parameters of
// type string; Name - which is the name of the bot, Description - which
// is the description of the bot and Industry - which is the related 
// industry of the bot.
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
	// infoLog.Println(string(pretty.Pretty(body)))
	if err := json.Unmarshal(body, &bot); err != nil {
		errorLog.Fatal("Cannot Unmarshal JSON response.")
	}
	chatID := uuid.New()
	bot.chatID = chatID.String()
	if resp.StatusCode == 200 {
		infoLog.Println("Bot created successfully")
	} else {
		errorLog.Fatal(string(pretty.Pretty(body)))
	}
	return bot
}

// UpdateBot() method to update contents on the bot. Used by other methods.
// User has the liberty to update manually too.
func (bot *Bot) UpdateBot() *Bot {

	if bot.Id == 0 {
		errorLog.Fatal("Cannot update a non existing bot")
	}
	infoLog.Print("Updating bot...")

	params := map[string]interface{}{
		"name":                 bot.Name,
		"description":          bot.Description,
		"industry":             bot.Industry,
		"visible_on_community": bot.VisibleOnCommunity,
		"intents":              bot.Intents,
		"flows":                bot.Flows,
	}

	jsonParams, err := json.Marshal(params)

	if err != nil {
		errorLog.Fatal(err)
	}
	url := fmt.Sprintf("%schatbot/%d", bot.baseURL, bot.Id)
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
	if resp.StatusCode != 200 {
		// infoLog.Println("Bot updated successfully")
		infoLog.Println(string(pretty.Pretty(body)))
	}
	return bot

}

// DeleteBot() method to delete the currently selected bot.
func (bot *Bot) DeleteBot() {
	if bot.Id == 0 {
		errorLog.Fatal("Cannot delete a non existing bot")
	}
	infoLog.Printf("Deleting bot with id: %d", bot.Id)
	url := fmt.Sprintf("%schatbot/%d", bot.baseURL, bot.Id)
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
	if resp.StatusCode != 200 {

		errorLog.Fatal(string(pretty.Pretty(body)))
	}
	bot.Id = 0
	bot.chatID = ""
	bot.Description = ""
	bot.Industry = ""
	bot.Name = ""

	for k := range bot.Intents {
		delete(bot.Intents, k)
	}
	for k := range bot.Flows {
		delete(bot.Flows, k)
	}

	infoLog.Println("Bot deleted successfully")
}
