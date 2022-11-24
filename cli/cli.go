package cli

import (
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/sarufi-io/sarufi-golang-sdk"
)

type (

	// Manager is the commandline tool to manages bots in sarufi platform
	Manager struct {
		mu       *sync.Mutex
		Logger   io.Writer
		Debug    bool
		http     *http.Client
		username string
		password string
		sarufi   *sarufi.Client
	}

	// Credentials
	Credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	// ManagerOption is the option for the manager
	ManagerOption func(*Manager)
)

// NewManager creates a new manager
func NewManager(credentials *Credentials, opts ...ManagerOption) *Manager {
	m := &Manager{
		mu:       &sync.Mutex{},
		Logger:   os.Stdout,
		Debug:    false,
		http:     http.DefaultClient,
		username: credentials.Username,
		password: credentials.Password,
	}
	for _, opt := range opts {
		opt(m)
	}

	c := sarufi.NewClient(
		sarufi.WithLogger(m.Logger),
		sarufi.WithDebug(m.Debug),
		sarufi.WithHttpClient(m.http),
	)
	m.sarufi = c
	return m
}

// SetCredentials sets the credentials for the manager
func (m *Manager) SetCredentials(credentials *Credentials) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.username = credentials.Username
	m.password = credentials.Password
}

// SetLogger sets the logger for the manager
func (m *Manager) SetLogger(logger io.Writer) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Logger = logger
	m.sarufi.SetLogger(logger)
}

// SetDebug sets the debug flag for the manager
func (m *Manager) SetDebug(debug bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Debug = debug
	m.sarufi.SetDebug(debug)
}

// SetHttpClient sets the logger for the manager
func (m *Manager) SetHttpClient(http *http.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.http = http
	m.sarufi.SetHttpClient(http)
}

// Run is the entrypoint for the manager
func (m *Manager) Run(args []string) error {
	return nil
}
