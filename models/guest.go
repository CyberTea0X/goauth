package models

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Guest struct {
	Id       int64  `json:"id,omitempty"`
	FullName string `json:"full_name,omitempty"`
}

// Calls rows.Next() and Scans the row into the guest struct
func (g *Guest) FromRow(rows *sql.Rows) (*Guest, error) {
	exists := rows.Next()
	if !exists {
		return g, errors.New(fmt.Sprintf("Can't scan guest from row"))
	}

	err := rows.Scan(&g.Id, &g.FullName)
	if err != nil {
		return g, errors.New(fmt.Sprintf("Can't scan guest from row: %s", err.Error()))
	}

	return g, nil
}

func RegisterGuest(fullname string, adress url.URL, client *http.Client) (*Guest, error) {
	//client := &http.Client{
	guest := new(Guest)
	guest.FullName = fullname
	requestBody, err := json.Marshal(guest)
	if err != nil {
		return nil, errors.Join(errors.New("Error marshalling guest struct to json"), err)
	}
	r, err := client.Post(adress.String(), "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errors.Join(errors.New("Error reading body of a new guest request response"), err)
	}
	if err := json.Unmarshal(body, guest); err != nil {
		return nil, errors.Join(errors.New("Error unmarshalling guest from request"), err)
	}

	return guest, nil
}
