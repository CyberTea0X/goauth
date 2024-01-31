package models

import (
	"io"
	"net/http"
)

type ClientMock struct {
	Response *http.Response
	Error    error
}

func (m *ClientMock) Post(string, string, io.Reader) (*http.Response, error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Response, nil
}

func (m *ClientMock) Get(url string) (resp *http.Response, err error) {
	if m.Error != nil {
		return nil, m.Error
	}
	return m.Response, nil
}
