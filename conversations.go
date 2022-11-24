package sarufi

import (
	"context"
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
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": fmt.Sprintf("Bearer %s", req.Token),
	}
	request, err := createRequest(ctx, method, reqURL, req, headers)
	if err != nil {
		return nil, fmt.Errorf("make conversation: %w", err)
	}
	resp, err := client.Do(request)
	defer resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("make conversation: %w", err)
	}
	var conversationResponse ConversationResponse
	err = parseResponse(resp, &conversationResponse)
	if err != nil {
		return nil, fmt.Errorf("make conversation: %w", err)
	}

	return &conversationResponse, nil
}
