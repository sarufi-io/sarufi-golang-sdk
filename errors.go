package sarufi

import "fmt"

type NotFoundError struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (nf *NotFoundError) Error() string {
	if nf != nil {
		if nf.Message != "" {
			return fmt.Sprintf("status code 404: %s", nf.Message)
		} else if nf.Detail != "" {
			return fmt.Sprintf("status code 404: %s", nf.Detail)
		}
	}
	return ""
}

type Unauthorized struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (ua *Unauthorized) Error() string {
	if ua != nil {
		if ua.Message != "" {
			return fmt.Sprintf("status code 401: %s", ua.Message)
		} else if ua.Detail != "" {
			return fmt.Sprintf("status code 401: %s", ua.Detail)
		}
	}
	return ""
}

type Detail struct {
	Loc     []string `json:"loc"`
	Message string   `json:"msg"`
	Type    string   `json:"type"`
}

type ConflictError struct {
	Detail Detail `json:"detail"`
}

func (c *ConflictError) Error() string {
	if c != nil {
		return fmt.Sprintf("status code 409: %s", c.Detail.Message)
	}
	return ""
}

type UnprocessableEntity struct {
	Detail Detail `json:"detail"`
}

func (ue *UnprocessableEntity) Error() string {
	if ue != nil {
		return fmt.Sprintf("status code 422: %s", ue.Detail.Message)
	}
	return ""
}
