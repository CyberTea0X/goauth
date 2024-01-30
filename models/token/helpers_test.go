package token

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const token = "sometoken"

func testExtractToken(t *testing.T, request *http.Request) (string, error) {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	var err error
	var ext string

	router.GET("/test", func(c *gin.Context) {
		ext, err = ExtractToken(c)
	})

	w := httptest.NewRecorder()

	router.ServeHTTP(w, request)

	return ext, err
}

func TestHeaderExtractSucceed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Add("Authorization", token)
	ext, err := testExtractToken(t, req)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token, ext)
}

func TestQueryExtractSucceed(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)

	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()

	ext, err := testExtractToken(t, req)

	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, token, ext)
}

func TestQueryExtractFailed(t *testing.T) {
	req, err := http.NewRequest("GET", "/test", nil)

	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	_, err = testExtractToken(t, req)

	if err == nil {
		t.Fatal("Token not specified, but no error was returned")
	}
}
