package sarufi

import (
	"log"
	"os"
	"time"
)

const BaseURL = "https://api.sarufi.io/"

// For Error and Information Loggings
var infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
var errorLog = log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime)

// A helper function to check if a file exists. It
// accepts a filename string and returns a bool.
func fileChecker(fileName string) bool {

	if len(fileName) > 50 {
		return false
	}
	_, error := os.Stat(fileName)

	// check if error is "file not exists"
	if !os.IsNotExist(error) {
		return true
	} else {
		return false
	}
}

func checkToken(token string) bool {
	if token != "" {
		return true
	}
	return false
}

type CustomTime struct {
	time.Time
}

func (m *CustomTime) UnmarshalJSON(data []byte) error {
	// Ignore null, like in the main JSON package.
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	// Fractional seconds are handled implicitly by Parse.
	tt, err := time.Parse(`"`+time.RFC3339+`"`, string(data))
	*m = CustomTime{tt}
	return err
}
