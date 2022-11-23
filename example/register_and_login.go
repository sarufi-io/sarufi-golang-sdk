package main

import (
	"encoding/json"
	"fmt"
	"github.com/sarufi-io/sarufi-golang-sdk"
)

func main() {

	//client := sarufi.NewClient(
	//	sarufi.WithHttpClient(http.DefaultClient),
	//	sarufi.WithLogger(os.Stdout),
	//	sarufi.WithDebug(true))
	//
	//// Register a new user
	//registerRequest := &sarufi.RegisterRequest{
	//	Username: "",
	//	Password: "",
	//}
	//
	//registerResponse, err := client.Register(context.Background(), registerRequest)
	//if err != nil {
	//	fmt.Printf("error: %s\n", err)
	//}
	//
	//fmt.Printf("login response: %+v\n", registerResponse)
	//
	//// Login
	//loginRequest := &sarufi.LoginRequest{
	//	Username: "",
	//	Password: "",
	//}
	//
	//loginResponse, err := client.Login(context.Background(), loginRequest)
	//if err != nil {
	//	fmt.Printf("error: %s\n", err)
	//}
	//
	//fmt.Printf("login response: %+v\n", loginResponse)

	/// Create a request like this
	//{
	//  "name": "go-sdk-bot",
	//  "description": "PUT DESCRIPTION HERE",
	//  "intents": {
	//      "version":["What is go version","Is go 2 out","Is this stable"],
	//      "maintainer" :["who wrote this crappy code","who is a maintainer","I want to contribute"]
	//  },
	//  "flows": {
	//      "version" :{
	//          "message" : ["the current version is 1.19","run go version in cmd","check it at https://go.dev"],
	//          "next_state": "end"
	//      },
	//      "maintainer":{
	//          "message":["see out github repository"],
	//          "next_state":"end"
	//      }
	//  },
	//  "industry": "general",
	//  "visible_on_community": true
	//}

	createReq := &sarufi.ChatbotCreateReq{
		Name:        "go-sdk-bot",
		Description: "PUT DESCRIPTION HERE",
		Intents: map[string][]string{
			"version":    {"What is go version", "Is go 2 out", "Is this stable"},
			"maintainer": {"who wrote this crappy code", "who is a maintainer", "I want to contribute"},
		},
		Flows: map[string]sarufi.Flow{
			"version": {
				Message:   []string{"the current version is 1.19", "run go version in cmd", "check it at https://go.dev"},
				NextState: "end",
			},
			"maintainer": {
				Message:   []string{"see out github repository"},
				NextState: "end",
			},
		},
		Industry:           "general",
		VisibleOnCommunity: true,
	}

	// marshal the request to json
	createReqJson, err := json.Marshal(createReq)
	if err != nil {
		fmt.Printf("error: %s\n", err)
	}

	// print the json
	fmt.Printf("create request: %s\n", string(createReqJson))

}
