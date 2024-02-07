package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

type User struct {
	Id       int64
	Username string
	Password string
	Email    string
	Roles    []string
}

type LoginServiceResponce struct {
	Id    *int64   `json:"id"`
	Roles []string `json:"roles"`
}

// Should return *User struct with all fields filled.
//
// returns an error if there are no roles or ID in the external service response
func LoginUser(client HTTPClient, adress url.URL, username string, password string, email string) (*User, error) {
	q := adress.Query()
	q.Set("username", username)
	q.Set("password", password)
	q.Set("email", email)
	adress.RawQuery = q.Encode()
	r, err := client.Get(adress.String())
	if err != nil {
		return nil, errors.Join(errors.New("HTTP error sending get request to login service"), err)
	}
	if r.StatusCode != 200 {
		extError := NewExternalServiceError(r.StatusCode)
		if err := ErrFromResponse(r); err != nil {
			extError.Msg = err.Error()
		}
		return nil, extError
	}
	rawBody, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Join(errors.New("error reading login responce body"), err)
	}
	res := new(LoginServiceResponce)
	if err := json.Unmarshal(rawBody, res); err != nil {
		return nil, errors.Join(errors.New("error unmarshalling login responce to User struct"), err)
	}
	if res.Id == nil {
		return nil, errors.New("no id specified in external service responce")
	}
	if res.Roles == nil {
		return nil, errors.New("no roles specified in external service responce")
	}
	user := new(User)
	user.Id = *res.Id
	user.Roles = res.Roles
	user.Username = username
	user.Password = password
	user.Email = email
	return user, nil
}
