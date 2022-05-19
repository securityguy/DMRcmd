/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package ha

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Config struct {
	Enabled bool   `json:"enabled"`
	Server  string `json:"server,omitempty"`
	Token   string `json:"token,omitempty"`
}

type HomeAssistant struct {
	Config Config
}

func New(c Config) HomeAssistant {
	var ha HomeAssistant
	ha.Config.Enabled = c.Enabled
	ha.Config.Server = c.Server
	ha.Config.Token = c.Token
	return ha
}

func (ha *HomeAssistant) Scene(s string) error {
	url := ha.Config.Server + "/api/services/scene/turn_on"
	entity := "scene." + s
	return ha.post(url, entity)
}

func (ha *HomeAssistant) Script(s string) error {
	url := ha.Config.Server + "/api/services/script/turn_on"
	entity := "script." + s
	return ha.post(url, entity)
}

func (ha *HomeAssistant) post(url string, entity string) error {

	// Create JSON message
	values := map[string]string{"entity_id": entity}
	jsonData, err := json.Marshal(values)
	if err != nil {
		return err
	}

	// Instantiate client
	client := &http.Client{Timeout: time.Second * 30}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	// Set headers
	req.Header.Set("Content-type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ha.Config.Token)

	// Execute request
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// Check response code
	if resp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("HTTP response code %d", resp.StatusCode))
	}
	return nil
}
