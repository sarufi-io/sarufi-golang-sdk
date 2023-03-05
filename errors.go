package sarufi

import "fmt"

type NotFoundError struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

func (nf *NotFoundError) Error() string {
	if nf != nil {
		if nf.Message != "" {
			return fmt.Sprintf("Not Found: %s", nf.Message)
		} else if nf.Detail != "" {
			return fmt.Sprintf("Not Found: %s", nf.Detail)
		}
	}
	return ""
}

type Unauthorized struct {
	Message string `json:"message"`
}

func (ua *Unauthorized) Error() string {
	if ua != nil {
		return fmt.Sprintf("Unauthorized: %s", ua.Message)
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
		return fmt.Sprintf("Conflict Response: %s", c.Detail.Message)
	}
	return ""
}

type UnprocessableEntity struct {
	Detail Detail `json:"detail"`
}

func (ue *UnprocessableEntity) Error() string {
	if ue != nil {
		return fmt.Sprintf("Unprocessable Entity: %s", ue.Detail.Message)
	}
	return ""
}
