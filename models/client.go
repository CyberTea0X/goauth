package models

import (
	"io"
	"net/http"
)

// Little interface so i can mock some requests
type HTTPClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
	Get(url string) (resp *http.Response, err error)
}
