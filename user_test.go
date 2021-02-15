package engage

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	c, _       = New(testKey, testSecretKey)
	testID     = "u1421"
	mockClient *mockHTTPClient
	body       map[string]interface{}
)

func init() {
	mockClient = &mockHTTPClient{}
	body = make(map[string]interface{})
	body["status"] = "ok"
	b, _ := json.Marshal(body)
	r := ioutil.NopCloser(bytes.NewReader(b))
	mockClient.DoFunc = func(req *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       r,
		}, nil
	}
}

func TestIdentifyUser(t *testing.T) {
	var (
		invalidIDData = map[string]interface{}{"email": "invalidEmail"}
	)
	_, err := c.User.Identify(nil)
	assert.NotNil(t, err)
	assert.Equal(t, errInvalidUserData, err)

	_, err = c.User.Identify(invalidIDData)
	assert.NotNil(t, err)
	assert.Equal(t, errInvalidOrMissingID, err)

	invalidIDData["id"] = "u141"
	_, err = c.User.Identify(invalidIDData)
	assert.NotNil(t, err)
	assert.Equal(t, errInvalidOrMissingEmail, err)

	data := map[string]interface{}{
		"id":    testID,
		"email": "test@engage.so",
	}

	// update client
	c.SetClient(mockClient)

	res, err := c.User.Identify(data)
	assert.Nil(t, err)
	assert.Equal(t, body, res)
}

func TestAddUserAttribute(t *testing.T) {
}
