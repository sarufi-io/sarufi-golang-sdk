package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func (bot *Bot) Initialize(username, password string) {
  bot.baseURL = "https://api.sarufi.io/"
  infoLog.Println("Getting Token...")
  url := bot.baseURL + "users/login"
  params := map[string]string{
    "username": username,
    "password": password,
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

  var result map[string]string
  json.Unmarshal([]byte(body), &result)

  if result["token"] != "" {
    bot.token = result["token"]
  } else {
    errorLog.Fatal(result["message"])
  }

}
