package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type (
	ChatbotRequestError struct {
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
)

// chatbotCreate creates a chatbot. It returns a Chatbot in case the request is successful
// returns content of the response body and an error if the request fails and the
// code is not 400. If the code is 400, it returns a ChatbotRequestError contents.
func chatbotCreate(ctx context.Context, client *http.Client, url, method string, request *ChatbotCreateReq) (*Chatbot, error) {
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
		var errResp ChatbotRequestError
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("chatbot create: %w", err)
		}
		return nil, fmt.Errorf("chatbot create: %s", errResp.Detail)
	} else if resp.StatusCode == http.StatusOK {
		var respBody Chatbot
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

// getChatBot gets a chatbot. It returns a Chatbot in case the request is successful
func getChatbot(ctx context.Context, client *http.Client, reqURL, method string, request *GetChatBotReq) (*Chatbot, error) {
	chatBotID := fmt.Sprintf("%d", request.ID)
	getChatbotURL, err := url.JoinPath(reqURL, chatBotID)
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, getChatbotURL, nil)
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("get chatbot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respBody Chatbot
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return nil, fmt.Errorf("get chatbot: %w", err)
		}
		return &respBody, nil
	} else if resp.StatusCode == 400 {
		var errResp ChatbotRequestError
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("get chatbot: %w", err)
		}
		return nil, fmt.Errorf("get chatbot: %s", errResp.Detail)
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("get chatbot: read response body: %w", err)
		}
		return nil, fmt.Errorf("get chatbot: %s", responseBody.String())
	}
}

func listChatBots(ctx context.Context, client *http.Client, reqURL, method, token string) ([]*Chatbot, error) {
	req, err := http.NewRequestWithContext(ctx, method, reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("list chatbots: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("list chatbots: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respBody []*Chatbot
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return nil, fmt.Errorf("list chatbots: %w", err)
		}
		return respBody, nil
	} else if resp.StatusCode == 400 {
		var errResp ChatbotRequestError
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("list chatbots: %w", err)
		}
		return nil, fmt.Errorf("list chatbots: %s", errResp.Detail)
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("list chatbots: read response body: %w", err)
		}
		return nil, fmt.Errorf("list chatbots: %s", responseBody.String())
	}
}

// updateChatbot updates a chatbot. It returns a Chatbot in case the request is successful
func updateChatbot(ctx context.Context, client *http.Client, reqURL, method string, request *UpdateChatbotReq) (*Chatbot, error) {
	chatBotID := fmt.Sprintf("%d", request.ID)
	updateChatbotURL, err := url.JoinPath(reqURL, chatBotID)
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	reqBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, updateChatbotURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Token))
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("update chatbot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		var respBody Chatbot
		err = json.NewDecoder(resp.Body).Decode(&respBody)
		if err != nil {
			return nil, fmt.Errorf("update chatbot: %w", err)
		}
		return &respBody, nil
	} else if resp.StatusCode == 400 {
		var errResp ChatbotRequestError
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return nil, fmt.Errorf("update chatbot: %w", err)
		}
		return nil, fmt.Errorf("update chatbot: %s", errResp.Detail)
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("update chatbot: read response body: %w", err)
		}
		return nil, fmt.Errorf("update chatbot: %s", responseBody.String())
	}
}

// deleteChatbot deletes a chatbot. It returns a Chatbot in case the request is successful
func deleteChatbot(ctx context.Context, client *http.Client, reqURL, method string, request *DeleteChatbotReq) error {
	chatBotID := fmt.Sprintf("%d", request.ID)
	deleteChatbotURL, err := url.JoinPath(reqURL, chatBotID)
	if err != nil {
		return fmt.Errorf("delete chatbot: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, method, deleteChatbotURL, nil)
	if err != nil {
		return fmt.Errorf("delete chatbot: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", request.Token))
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("delete chatbot: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode == 400 {
		var errResp ChatbotRequestError
		err = json.NewDecoder(resp.Body).Decode(&errResp)
		if err != nil {
			return fmt.Errorf("delete chatbot: %w", err)
		}
		return fmt.Errorf("delete chatbot: %s", errResp.Detail)
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return fmt.Errorf("delete chatbot: read response body: %w", err)
		}
		return fmt.Errorf("delete chatbot: %s", responseBody.String())
	}
}
