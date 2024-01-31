package models

import (
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

type User struct {
	Id       int64  `json:"id"`
	Username string `json:"username"` // not null
	Password string `json:"password"` // not null
	Email    string `json:"email"`    // not null unique
	Role     string `json:"role"`     // not null
}

func LoginUser(client HTTPClient, adress url.URL, username string, password string, email string) (*User, error) {
	user := new(User)
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
		return nil, errors.Join(errors.New("Error reading login responce body"), err)
	}
	if err := json.Unmarshal(rawBody, user); err != nil {
		return nil, errors.Join(errors.New("Error unmarshalling login responce to User struct"), err)
	}
	user.Username = username
	user.Password = password
	user.Email = email
	return user, nil
}

func (u *User) PrepareGive() {
	u.Password = ""
	u.Email = ""
}
