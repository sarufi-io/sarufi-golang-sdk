package sarufi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// parseResponse takes a http.Response and a v which is a struct to which the body
// will be unmarshalled add returns an error
func parseResponse(response *http.Response, v any) error {
	statusCode := response.StatusCode
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
