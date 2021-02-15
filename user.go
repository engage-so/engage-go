package engage

import (
	"errors"
	"fmt"
	"regexp"
)

// UserResource handles all resource operations for Users
type UserResource resource

var (
	errInvalidUserData       = errors.New("You need to pass an object with at least and id and email")
	errInvalidOrMissingID    = errors.New("ID is missing")
	errInvalidOrMissingEmail = errors.New("Email is missing or invalid")
	errNoAttributeUser       = errors.New("User id missing")
	errNoAttributeData       = errors.New("Attributes data is missing")

	// Regexp for validating email
	re = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

	// Allowed attributes
	allowed = []string{"id", "email", "device_token", "device_platform", "number", "created_at", "first_name", "last_name"}

	// Non meta attributes
	nonMeta = []string{"email", "device_token", "device_platform", "number", "created_at", "first_name", "last_name"}
)

// Identify Engage user
func (r *UserResource) Identify(data map[string]interface{}) (identifyResp map[string]interface{}, err error) { // validate data
	if data == nil {
		return nil, errInvalidUserData
	}

	if _, ok := data["id"]; !ok {
		return nil, errInvalidOrMissingID
	}

	if _, ok := data["email"]; !ok || !re.MatchString(data["email"].(string)) {
		return nil, errInvalidOrMissingEmail
	}

	params := make(map[string]interface{})

	for k, v := range data {
		if indexOf(k, allowed) != -1 {
			params[k] = v
		}
	}
	endpoint := fmt.Sprintf("/users/%s", data["id"])
	resp, err := r.client.putRequest(endpoint, params)
	if err != nil {
		return
	}

	err = resp.ParseJSON(&identifyResp)
	return
}

// AddAttribute add attributes to users for segmentation
func (r *UserResource) AddAttribute(userid string, data map[string]interface{}) (attributeResp map[string]interface{}, err error) {
	if userid == "" {
		return nil, errNoAttributeUser
	}
	if data == nil {
		return nil, errNoAttributeData
	}

	params := map[string]interface{}{
		"meta": map[string]interface{}{},
	}

	for k, v := range data {
		if indexOf(k, nonMeta) != -1 {
			params[k] = v
		} else {
			params["meta"].(map[string]interface{})[k] = v
		}
	}

	endpoint := fmt.Sprintf("/users/%s", userid)
	resp, err := r.client.putRequest(endpoint, params)
	if err != nil {
		return
	}

	err = resp.ParseJSON(&attributeResp)

	return
}

// Track enables events and user actions to be tracked
func (r *UserResource) Track(userid string, data interface{}) (trackResponse map[string]interface{}, err error) {
	var payload map[string]interface{}

	if userid == "" {
		return nil, errNoAttributeUser
	}

	if data == nil {
		return nil, errNoAttributeData
	}

	// validate and prepare payload
	switch data.(type) {
	case string:
		payload = map[string]interface{}{
			"event": data,
			"value": true,
		}
	case map[string]interface{}:
		payload = data.(map[string]interface{})
	default:
		return nil, errNoAttributeData
	}

	endpoint := fmt.Sprintf("/users/%s/events", userid)
	resp, err := r.client.putRequest(endpoint, payload)
	if err != nil {
		return
	}

	err = resp.ParseJSON(&trackResponse)

	return
}

func indexOf(word string, data []string) int {
	for k, v := range data {
		if word == v {
			return k
		}
	}
	return -1
}
