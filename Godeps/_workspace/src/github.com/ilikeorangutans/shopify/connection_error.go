package shopify

import (
	"fmt"
	"net/http"
)

func FromResponse(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 400 {
		return nil
	}

	return ConnectionError{
		Method:     resp.Request.Method,
		URL:        resp.Request.URL.String(),
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
	}
}

type ConnectionError struct {
	Method     string
	URL        string
	StatusCode int
	Status     string
}

func (ce ConnectionError) Error() string {
	return fmt.Sprintf("%s %s -> %d %s", ce.Method, ce.URL, ce.StatusCode, ce.Status)
}
