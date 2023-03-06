package sarufi

import (
	"fmt"
	"io"
	"net/http"
)

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
