package pixiv

import "errors"

// APIError returned by the api
type APIError struct {
	ErrorDetails struct {
		UserMessage string `json:"user_message,omitempty"`
		Message     string `json:"message,omitempty"`
		Reason      string `json:"reason,omitempty"`
	} `json:"error,omitempty"`
}

func (e *APIError) Error() string {
	return e.ErrorDetails.Message
}

// HasError checks if a message is set
func (e *APIError) HasError() bool {
	return e.Error() != ""
}

// The errors
var (
	ErrAuthentication = errors.New("failed to authenticate")
	ErrMissingToken   = errors.New("missing oauth token")
)
