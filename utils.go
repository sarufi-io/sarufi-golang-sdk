package sarufi

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const baseURL = "https://api.sarufi.io/"

var token = ""

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

// A helper function to make requests easier
func makeRequest(method, url string, data io.Reader) (int, []byte, error) {
	req, err := http.NewRequest(method, url, data)

	if err != nil {
		return 0, nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	bearer := fmt.Sprintf("Bearer %s", token)
	req.Header.Set("Authorization", bearer)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return 0, nil, err
	}

	return resp.StatusCode, body, nil
}

// ToDo;
// Make custom way of parsing time
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
