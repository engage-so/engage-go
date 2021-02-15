package engage

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testKey       = "mytestkey"
	testSecretKey = "mytestsecretkey"
)

var defaultClient = &http.Client{Timeout: defaultTimeout}

// mockHTTPClient
type mockHTTPClient struct {
	DoFunc func(*http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestEngage(t *testing.T) {
	c, err := New(testSecretKey, testSecretKey)
	assert.Nil(t, err, "expected client creation to yield nil errors")
	c.SetClient(defaultClient)
	assert.Equal(t, c.client, defaultClient, "expected set client to match input")
}

func TestCannotInitWithoutKeys(t *testing.T) {
	_, err := New("", "")
	assert.NotNil(t, err)
	assert.Equal(t, errNoAPIKey, err)
}

func TestHTTPClient(t *testing.T) {
	mockClient := &mockHTTPClient{}
	body := `{"status": "ok"}`
	r := ioutil.NopCloser(bytes.NewReader([]byte(body)))

	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
	c, _ := New(testSecretKey, testSecretKey)
	c.SetClient(mockClient)
	payload := map[string]string{"event": "sample"}
	res, err := c.postRequest("/users/u141", payload)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)

	res, err = c.putRequest("/users/u141", payload)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.Code)
}
