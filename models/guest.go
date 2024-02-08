package models

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/url"
)

type Guest struct {
	Id   int64  `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// model.ExternalServiceError error returned if error is parsed from external service response
func RegisterGuest(fullname string, adress url.URL, client HTTPClient) (*Guest, error) {
	guest := new(Guest)
	guest.Name = fullname
	requestBody, err := json.Marshal(guest)
	if err != nil {
		return nil, errors.Join(errors.New("Error marshalling guest struct to json"), err)
	}
	r, err := client.Post(adress.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errors.Join(errors.New("HTTP error sending register post request"), err)
	}
	defer r.Body.Close()
	if r.StatusCode != 200 {
		extError := NewExternalServiceError(r.StatusCode)
		if err := ErrFromResponse(r); err != nil {
			extError.Msg = err.Error()
		}
		return nil, extError
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Error reading body of a new guest request response"), err)
	}
	if err := json.Unmarshal(body, guest); err != nil {
		return nil, errors.Join(errors.New("Error unmarshalling guest from request"), err)
	}

	return guest, nil
}
