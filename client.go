package sarufi

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var (
	_ Service             = (*Client)(nil)
	_ BotService          = (*Client)(nil)
	_ ConversationService = (*Client)(nil)
)

type (
	// Client is the sarufi client
	Client struct {
		http   *http.Client
		logger io.Writer
		debug  bool
	}

	ClientOption func(*Client)
)

// SetLogger sets the logger for the client
func (c *Client) SetLogger(logger io.Writer) {
	c.logger = logger
}

// SetDebug sets the debug flag for the client
func (c *Client) SetDebug(debug bool) {
	c.debug = debug
}

// SetHttpClient sets the logger for the client
func (c *Client) SetHttpClient(http *http.Client) {
	c.http = http
}

// WithHttpClient sets the http client for the client
func WithHttpClient(http *http.Client) ClientOption {
	return func(c *Client) {
		c.http = http
	}
}

// WithLogger sets the logger for the client
func WithLogger(logger io.Writer) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithDebug sets the debug flag for the client
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		http:   http.DefaultClient,
		logger: os.Stdout,
		debug:  false,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *Client) Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error) {
	loginURL, err := url.JoinPath(BaseURL, UsersLoginEndpoint)
	if err != nil {
		return nil, fmt.Errorf("login: join url: %w", err)
	}
	return login(ctx, c.http, loginURL, http.MethodPost, request)
}

func (c *Client) Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error) {
	registerURL, err := url.JoinPath(BaseURL, UsersRegisterEndpoint)
	if err != nil {
		return nil, fmt.Errorf("register: join url: %w", err)
	}
	return register(ctx, c.http, registerURL, http.MethodPost, request)
}

func (c *Client) MakeConversation(ctx context.Context, req *ConversationRequest) (*ConversationResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (c *Client) CreateBot(ctx context.Context, req *ChatbotCreateReq) (*Chatbot, error) {
	// TODO implement me
	panic("implement me")
}

func (c *Client) GetBot(ctx context.Context, req *GetChatBotReq) (*Chatbot, error) {
	// TODO implement me
	panic("implement me")
}

func (c *Client) DeleteBot(ctx context.Context, req *DeleteChatbotReq) error {
	// TODO implement me
	panic("implement me")
}

func (c *Client) UpdateBot(ctx context.Context, req *UpdateChatbotReq) (*Chatbot, error) {
	// TODO implement me
	panic("implement me")
}

func (c *Client) ListBots(ctx context.Context, token string) ([]*Chatbot, error) {
	// TODO implement me
	panic("implement me")
}
