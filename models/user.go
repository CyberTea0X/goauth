package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int64
	Username string
	Password string
	Email    string
	Roles    []string
}

type LoginResponce struct {
	Id           *int64   `json:"id" validate:"required"`
	Roles        []string `json:"roles" validate:"required"`
	PasswordHash string   `json:"password" validate:"required"`
	HashAlg      string   `json:"alg" validate:"required"`
}

type LoginRequest struct {
	Username string
	Email    string
}

// Should return *User struct with all fields filled.
//
// returns an error if there are no roles or ID in the external service response
func LoginUser(client HTTPClient, adress url.URL, username string, password string, email string) (*User, error) {
	q := adress.Query()
	req := LoginRequest{Username: username, Email: password}
	q.Set("login", req.Username)
	q.Set("email", req.Email)
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
	res := new(LoginResponce)
	if err := json.Unmarshal(rawBody, res); err != nil {
		return nil, errors.Join(errors.New("error unmarshalling login responce to User struct"), err)
	}
	validate := validator.New(validator.WithRequiredStructEnabled())
	err = validate.Struct(res)
	if err != nil {
		return nil, errors.Join(errors.New("invalid json responce from external login service"), err)
	}
	if err := ValidatePassword(password, res.PasswordHash, res.HashAlg); err != nil {
		return nil, err
	}
	user := new(User)
	user.Id = *res.Id
	user.Roles = res.Roles
	user.Username = username
	user.Password = password
	user.Email = email
	return user, nil
}

func ValidatePassword(password string, passwordHash, alg string) error {
	if alg != "bcrypt" {
		return fmt.Errorf("bad hashing algorithm from external login service. Expected bcrypt, got %s", alg)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password)); err != nil {
		errors.Join(errors.New("invalid password"), err)
	}
	return nil
}
