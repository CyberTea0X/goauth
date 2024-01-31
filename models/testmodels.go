package models

import (
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

type ClientMock struct {
	Engine   *gin.Engine
	Recorder *httptest.ResponseRecorder
}

func NewClientMock() *ClientMock {
	return &ClientMock{
		Engine:   gin.Default(),
		Recorder: httptest.NewRecorder(),
	}
}

func (m *ClientMock) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	r, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	r.Header.Set("Content-Type", contentType)
	m.Engine.ServeHTTP(m.Recorder, r)
	return m.Recorder.Result(), nil
}

func (m *ClientMock) Get(url string) (resp *http.Response, err error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	m.Engine.ServeHTTP(m.Recorder, r)
	return m.Recorder.Result(), nil
}
