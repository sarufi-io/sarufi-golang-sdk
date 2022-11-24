package sarufi

import (
	"context"
)

const (
	BaseURL               = "https://api.sarufi.io"
	UsersLoginEndpoint    = "/users/login"
	UsersRegisterEndpoint = "/users/register"
	ChatbotEndpoint       = "/chatbot"
	ListChatbotsEndpoint  = "/chatbots"
)

type (

	// Service is a sarufi api interface ..
	Service interface {
		Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
		Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error)
		// ChatbotCreate(ctx context.Context, request *ChatbotCreateReq) (*Chatbot, error)
	}
)
