package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

const (
	EnvSarufiPlatformUsername = "SARUFI_PLATFORM_USERNAME"
	EnvSarufiPlatformPassword = "SARUFI_PLATFORM_PASSWORD"
	DefSarufiPlatformUsername = ""
	DefSarufiPlatformPassword = ""
)

type (

	// Credentials ..
	Credentials struct {
		Username string
		Password string
	}
	// dotEnvReader implements the Reader interface
	// It reads the environment variables from a .env file
	// and stores them in a map
	dotEnvReader struct {
		rwm      *sync.RWMutex
		filename string
		env      map[string]string
	}
)

// LoadCredentials
func LoadCredentialsFromEnv(reader Reader) *Credentials {
	return &Credentials{
		Username: reader.String(EnvSarufiPlatformUsername, DefSarufiPlatformUsername),
		Password: reader.String(EnvSarufiPlatformPassword, DefSarufiPlatformPassword),
	}
}

// NewDotEnvReader returns a new instance of the Reader interface
func NewDotEnvReader(filename string) (Reader, error) {
	m, err := godotenv.Read(filename)
	if err != nil {
		return nil, fmt.Errorf("error reading env file(%s): %v", filename, err)
	}
	return &dotEnvReader{
		rwm:      &sync.RWMutex{},
		filename: filename,
		env:      m,
	}, nil
}

func (reader *dotEnvReader) Add(key string, value string) error {
	reader.rwm.Lock()
	defer reader.rwm.Unlock()
	if _, ok := reader.env[key]; ok {
		return fmt.Errorf("key(%s) already exists", key)
	}
	reader.env[key] = value
	// update the file
	return godotenv.Write(reader.env, reader.filename)
}

func (reader *dotEnvReader) Update(key string, value string) error {
	reader.rwm.Lock()
	defer reader.rwm.Unlock()
	if _, ok := reader.env[key]; !ok {
		return fmt.Errorf("key(%s) does not exist", key)
	}
	reader.env[key] = value

	// update the file
	return godotenv.Write(reader.env, reader.filename)
}

// Delete ...
func (reader *dotEnvReader) Delete(key string) error {
	reader.rwm.Lock()
	defer reader.rwm.Unlock()
	delete(reader.env, key)
	return nil
}

func (reader *dotEnvReader) String(key string, defaultValue string) string {
	reader.rwm.RLock()
	defer reader.rwm.RUnlock()
	if v, ok := reader.env[key]; ok {
		return v
	}
	return defaultValue
}

func (reader *dotEnvReader) Int(key string, defaultValue int64) int64 {
	reader.rwm.RLock()
	defer reader.rwm.RUnlock()
	if v, ok := reader.env[key]; ok {
		// convert the string to int
		iv, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return defaultValue
		}
		return iv
	}
	return defaultValue
}

func (reader *dotEnvReader) Bool(key string, defaultValue bool) bool {
	reader.rwm.RLock()
	defer reader.rwm.RUnlock()
	if v, ok := reader.env[key]; ok {
		return v == "true" || v == "1"
	}
	return defaultValue
}

func (reader *dotEnvReader) Float(key string, defaultValue float64) float64 {
	reader.rwm.RLock()
	defer reader.rwm.RUnlock()
	if v, ok := reader.env[key]; ok {
		// convert the string to int
		iv, err := strconv.ParseFloat(v, 64)
		if err != nil {
			return defaultValue
		}
		return iv
	}
	return defaultValue
}

// DefaultDotEnvReader returns a new instance of the Reader interface
func DefaultDotEnvReader() (Reader, error) {
	ex, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("config reader: error getting executable path: %v", err)
	}
	exPath := filepath.Dir(ex)
	configPath := filepath.Join(exPath, "sarufi.env")
	reader, err := NewDotEnvReader(configPath)
	if err != nil {
		return nil, fmt.Errorf("config reader: error creating default reader: %v", err)
	}

	return reader, nil
}
