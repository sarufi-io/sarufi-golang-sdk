package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	GeneralConversation  ConversationType = "general"
	WhatsappConversation ConversationType = "whatsapp"
)

type (
	ConversationType    string
	ConversationRequest struct {
		Token       string           `json:"-"`
		Type        ConversationType `json:"-"`
		ChatID      string           `json:"chat_id"`
		BotID       string           `json:"bot_id"`
		Message     string           `json:"message"`
		MessageType string           `json:"message_type"`
		Language    string           `json:"language"`
	}

	ConversationResponse struct {
		Message   []string          `json:"message"`
		Memory    map[string]string `json:"memory"`
		NextState string            `json:"next_state"`
	}
)

// makeConversationRequest makes a request to the Sarufi API to get a response
// from the bot
func makeConversation(ctx context.Context, client *http.Client, reqURL, method string, req *ConversationRequest) (*ConversationResponse, error) {
	if req.Type == WhatsappConversation {
		reqURL = reqURL + "/whatsapp"
	}
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("make conversation: %w", err)
	}
	request, err := http.NewRequestWithContext(ctx, method, reqURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("make conversation: %w", err)
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+req.Token)
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("make conversation: %w", err)
	}

	// If Ok return the ConversationResponse, if 422 return a validation error
	// if 400 return a RequestError, if anything else return the response body
	// as an error
	if resp.StatusCode == http.StatusOK {
		var conversationResponse ConversationResponse
		if err := json.NewDecoder(resp.Body).Decode(&conversationResponse); err != nil {
			return nil, fmt.Errorf("make conversation: %w", err)
		}
		return &conversationResponse, nil
	} else if resp.StatusCode == 400 {
		var requestError RequestError
		if err := json.NewDecoder(resp.Body).Decode(&requestError); err != nil {
			return nil, fmt.Errorf("make conversation: %w", err)
		}
		return nil, fmt.Errorf("make conversation: %s", requestError.Detail)
	} else if resp.StatusCode == http.StatusUnprocessableEntity {
		var validationError ValidationError
		if err := json.NewDecoder(resp.Body).Decode(&validationError); err != nil {
			return nil, fmt.Errorf("make conversation: %w", err)
		}
		return nil, fmt.Errorf("make conversation: %w", &validationError)
	} else {
		var responseBody bytes.Buffer
		if _, err := responseBody.ReadFrom(resp.Body); err != nil {
			return nil, fmt.Errorf("make conversation: %w", err)
		}
		return nil, fmt.Errorf("make conversation: %s", responseBody.String())
	}
}
