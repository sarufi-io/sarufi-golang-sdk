package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

type (
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

	UpdateChatbotReq struct {
		ID                 int64               `json:"-"`
		Token              string              `json:"-"`
		Name               string              `json:"name"`
		Description        string              `json:"description"`
		Intents            map[string][]string `json:"intents"`
		Flows              map[string]Flow     `json:"flows"`
		Industry           string              `json:"industry"`
		VisibleOnCommunity bool                `json:"visible_on_community"`
	}

	// Chatbot this is a response sent after creating a chatbot
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
	Chatbot struct {
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

	// GetChatBotReq this is a request sent to get a chatbot
	GetChatBotReq struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}

	// DeleteChatbotReq this is a request sent to delete a chatbot
	DeleteChatbotReq struct {
		ID    int    `json:"id"`
		Token string `json:"token"`
	}

	BotService interface {
		CreateBot(ctx context.Context, req *ChatbotCreateReq) (*Chatbot, error)
		GetBot(ctx context.Context, req *GetChatBotReq) (*Chatbot, error)
		DeleteBot(ctx context.Context, req *DeleteChatbotReq) error
		UpdateBot(ctx context.Context, req *UpdateChatbotReq) (*Chatbot, error)
		ListBots(ctx context.Context, token string) ([]*Chatbot, error)
	}
)

// chatbotCreate creates a chatbot. It returns a Chatbot in case the request is successful
// returns content of the response body and an error if the request fails and the
// code is not 400. If the code is 400, it returns a RequestError contents.
func chatbotCreate(ctx context.Context, client *http.Client, url, method string, request *ChatbotCreateReq) (*Chatbot, error) {
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("chatbot create: %w", err)
	}
	req, err := createRequest(ctx, method, url, bytes.NewBuffer(reqBody),
		map[string]string{
			"Content-Type": "application/json",
		})
	if err != nil {
		return nil, fmt.Errorf("chatbot create: %w", err)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("chatbot create: %w", err)
	}
	var response Chatbot
	err = parseResponse(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("create chatbot: %w", err)
	}
	return &response, nil
}

// getChatBot gets a chatbot. It returns a Chatbot in case the request is successful
func getChatbot(ctx context.Context, client *http.Client, reqURL, method string, request *GetChatBotReq) (*Chatbot, error) {
	chatBotID := fmt.Sprintf("%d", request.ID)
	getChatbotURL, err := url.JoinPath(reqURL, chatBotID)
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", request.Token)}
	req, err := createRequest(ctx, method, getChatbotURL, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}
	var response Chatbot
	err = parseResponse(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}
	return &response, nil
}

func listChatBots(ctx context.Context, client *http.Client, reqURL, method, token string) ([]*Chatbot, error) {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", token)}
	req, err := createRequest(ctx, method, reqURL, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("list chatbots: %w", err)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("list chatbots: %w", err)
	}

	var response []*Chatbot
	err = parseResponse(resp, &response)
	if err != nil {
		return nil, fmt.Errorf("list chatbots: %w", err)
	}
	return response, nil
}

// updateChatbot updates a chatbot. It returns a Chatbot in case the request is successful
func updateChatbot(ctx context.Context, client *http.Client, reqURL, method string, request *UpdateChatbotReq) (*Chatbot, error) {
	chatBotID := fmt.Sprintf("%d", request.ID)
	updateChatbotURL, err := url.JoinPath(reqURL, chatBotID)
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", request.Token)}
	req, err := createRequest(ctx, method, updateChatbotURL, request, headers)
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	var respBody Chatbot
	err = parseResponse(resp, &respBody)
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	return &respBody, nil
}

// deleteChatbot deletes a chatbot. It returns a Chatbot in case the request is successful
func deleteChatbot(ctx context.Context, client *http.Client, reqURL, method string, request *DeleteChatbotReq) error {
	chatBotID := fmt.Sprintf("%d", request.ID)
	deleteChatbotURL, err := url.JoinPath(reqURL, chatBotID)
	if err != nil {
		return fmt.Errorf("delete chatbot: %w", err)
	}
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", request.Token)}
	req, err := createRequest(ctx, method, deleteChatbotURL, nil, headers)
	if err != nil {
		return fmt.Errorf("delete chatbot: %w", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delete chatbot: %w", err)
	}

	if resp != nil {
		defer resp.Body.Close()
	}
	return parseResponse(resp, nil)
}

// LoadConversationFlow loads a conversation flows from a json file that will be used to train the bot model
// It returns a []*Flow in case the request is successful. The file content should be in
// JSON format.
func LoadConversationFlow(path string) ([]*Flow, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load conversation flow: %w", err)
	}

	var flows []*Flow
	err = json.Unmarshal(f, &flows)
	if err != nil {
		return nil, fmt.Errorf("load conversation flow: %w", err)
	}

	return flows, nil
}

// LoadIntents loads intents from a file, It is exepected that the file contents are
// in JSON format.
func LoadIntents(path string) (Intents, error) {
	// read all the contents of the file into a byte slice
	// then unmarshal the byte slice into a Intents struct
	// and return the Intents struct
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load intents: %w", err)
	}
	var intents Intents
	err = json.Unmarshal(dataBytes, &intents)
	if err != nil {
		return nil, fmt.Errorf("load intents: %w", err)
	}
	return intents, nil
}
