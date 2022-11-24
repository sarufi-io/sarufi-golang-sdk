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
