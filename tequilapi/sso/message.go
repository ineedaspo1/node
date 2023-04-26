package sso

import (
	"encoding/json"
)

type MystnodesMessage struct {
	CodeChallenge string `json:"code_challenge"`
	Identity      string `json:"identity"`
	RedirectURL   string `json:"redirect_url"`
}

func (msg MystnodesMessage) json() ([]byte, error) {
	payload, err := json.Marshal(msg)
	if err != nil {
		return []byte{}, err
	}
	return payload, nil
}
