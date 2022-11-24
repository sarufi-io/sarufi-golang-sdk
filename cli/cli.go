package cli

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/sarufi-io/sarufi-golang-sdk"
	"github.com/urfave/cli/v2"
)

type (
	Commander struct {
	}

	// Manager is the commandline tool to manages bots in sarufi platform
	Manager struct {
		mu       *sync.Mutex
		Logger   io.Writer
		Debug    bool
		http     *http.Client
		username string
		password string
		sarufi   *sarufi.Client
		cli      *cli.App
	}

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
		cli: &cli.App{
			Name:  "sarufi",
			Usage: "sarufi is a commandline tool to manage bots in sarufi paltform",
			Action: func(*cli.Context) error {
				fmt.Println("boom! I say!")
				return nil
			},
		},
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

// LoginCommand ...
func LoginCommand(f sarufi.LoginFunc) *cli.Command {
	return &cli.Command{}
}

// RegisterCommand returns a cli.Command that is executed during registration
// of a new user
// sarufi register --username="johndoe" --password="johndoepassword"
func RegisterCommand(f sarufi.RegisterFunc) *cli.Command {
	cmnd := &cli.Command{
		Name:  "register",
		Usage: "register a new user to sarufi platform",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "username",
				Value: "",
				Usage: "username",
			},
			&cli.StringFlag{
				Name:  "password",
				Value: "",
				Usage: "password",
			},
		},
		Action: func(cCtx *cli.Context) error {
			// retrieve username and password

			fmt.Println("added task: ", cCtx.Args().First())
			return nil
		},
	}

	return cmnd
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
