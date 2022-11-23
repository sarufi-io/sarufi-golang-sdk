package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	BaseURL               = "https://api.sarufi.io"
	UsersLoginEndpoint    = "/users/login"
	UsersRegisterEndpoint = "/users/register"
	ChatbotEndpoint       = "/chatbot"
	ListChatbotsEndpoint  = "/chatbots"
)

var _ Service = (*Client)(nil)

type (
	// Client is the sarufi client
	Client struct {
		http   *http.Client
		logger io.Writer
		debug  bool
	}

	ClientOption func(*Client)

	// Service is a sarufi api interface ..
	Service interface {
		Login(ctx context.Context, request *LoginRequest) (*LoginResponse, error)
		Register(ctx context.Context, request *RegisterRequest) (*RegisterResponse, error)
		ChatbotCreate(ctx context.Context, request *ChatbotCreateReq) (*Chatbot, error)
	}
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

func (c *Client) ChatbotCreate(ctx context.Context, request *ChatbotCreateReq) (*Chatbot, error) {
	//TODO implement me
	panic("implement me")
}

// parseResponse takes a http.Response and a v which is a struct to which the body
// will be unmarshalled add returns an error
func parseResponse(response *http.Response, v any) error {
	statusCode := response.StatusCode
	defer response.Body.Close()
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("parse response: %w", err)
	}

	if statusCode == 200 || statusCode == 201 {
		decodeErr := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(v)
		if decodeErr != nil {
			return fmt.Errorf("parse response: %w", decodeErr)
		}
	} else if statusCode == 422 {
		var validationError ValidationError
		if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&validationError); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
		return fmt.Errorf("parse response: %w", &validationError)
	} else if statusCode == 400 {
		var requestError RequestError
		if err := json.NewDecoder(bytes.NewBuffer(bodyBytes)).Decode(&requestError); err != nil {
			return fmt.Errorf("parse response: %w", err)
		}
		return fmt.Errorf("parse response: %w", &requestError)
	} else {
		return fmt.Errorf("parse response: %s", string(bodyBytes))
	}

	return nil
}

// createRequest creates a http request with the given method, url and body and headers
func createRequest(ctx context.Context, method, url string, body any, headers map[string]string) (*http.Request, error) {
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}
	request, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	for key, value := range headers {
		request.Header.Set(key, value)
	}
	return request, nil
}
