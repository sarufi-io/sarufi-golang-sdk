package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type (
	ChatbotCreateError struct {
		Detail string `json:"detail"`
	}
	Intents map[string][]string
	Flow    struct {
		Message   []string `json:"message"`
		NextState string   `json:"next_state"`
	}
	// ChatbotCreateReq this is a request sent to create a chatbot
	// it is submitted as a json object like this:
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
	// Flows is simply a map of Flow
	// Intents is a map of string array
	ChatbotCreateReq struct {
		Name               string              `json:"name"`
		Description        string              `json:"description"`
		Intents            map[string][]string `json:"intents"`
		Flows              map[string]Flow     `json:"flows"`
		Industry           string              `json:"industry"`
		VisibleOnCommunity bool                `json:"visible_on_community"`
	}

	// ChatbotCreateResp this is a response sent after creating a chatbot
	// it is a json object like this:
	//{
	//    "intents": {
	//        "version": [
	//            "What is go version",
	//            "Is go 2 out",
	//            "Is this stable"
	//        ],
	//        "maintainer": [
	//            "who wrote this crappy code",
	//            "who is a maintainer",
	//            "I want to contribute"
	//        ]
	//    },
	//    "user_id": 74,
	//    "description": "PUT DESCRIPTION HERE",
	//    "industry": "general",
	//    "created_at": "2022-11-23T16:44:38.633418",
	//    "name": "go-sdk-bot",
	//    "flows": {
	//        "version": {
	//            "message": [
	//                "the current version is 1.19",
	//                "run go version in cmd",
	//                "check it at https://go.dev"
	//            ],
	//            "next_state": "end"
	//        },
	//        "maintainer": {
	//            "message": [
	//                "see out github repository"
	//            ],
	//            "next_state": "end"
	//        }
	//    },
	//    "id": 139,
	//    "model_name": "models/42faa73dcef83f9c89f4c8e2c45aa015.pkl",
	//    "visible_on_community": true,
	//    "updated_at": "2022-11-23T16:44:38.633430"
	ChatbotCreateResp struct {
		Intents            map[string][]string `json:"intents"`
		UserID             int                 `json:"user_id"`
		Description        string              `json:"description"`
		Industry           string              `json:"industry"`
		CreatedAt          string              `json:"created_at"`
		Name               string              `json:"name"`
		Flows              map[string]Flow     `json:"flows"`
		ID                 int                 `json:"id"`
		ModelName          string              `json:"model_name"`
		VisibleOnCommunity bool                `json:"visible_on_community"`
		UpdatedAt          string              `json:"updated_at"`
	}
)

// chatbotCreate creates a chatbot. It returns a ChatbotCreateResp in case the request is successful
// returns content of the response body and an error if the request fails and the
// code is not 400. If the code is 400, it returns a ChatbotCreateError contents.
func chatbotCreate(ctx context.Context, client *http.Client, url, method string, request *ChatbotCreateReq) (*ChatbotCreateResp, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("chatbot create: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("chatbot create: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("chatbot create: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusBadRequest {
		var errResp ChatbotCreateError
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("chatbot create: %w", err)
		}
		return nil, fmt.Errorf("chatbot create: %w", errResp)
	} else if resp.StatusCode == http.StatusOK {
		var respBody ChatbotCreateResp
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return nil, fmt.Errorf("chatbot create: %w", err)
		}
		return &respBody, nil
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("register: read response body: %w", err)
		}
		return nil, fmt.Errorf("chatbot create: %s", responseBody.String())
	}
}
