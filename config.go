/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
)

// Structure to hold configuration information
type configData struct {
	Debug   bool           `json:"debug"`
	Listen  string         `json:"listen"`
	Minimum uint32         `json:"minimum,omitempty"`
	Clients []configClient `json:"clients"`
	Events  []configEvent  `json:"events"`
}

type configClient struct {
	ID       uint32 `json:"id"`
	Password string `json:"password"`
}

type configEvent struct {
	Enabled   bool              `json:"enabled"`
	Name      string            `json:"name,omitempty"`
	SRC       uint32            `json:"src,omitempty"`
	DST       uint32            `json:"dst,omitempty"`
	Client    uint32            `json:"client,omitempty"`
	IP        string            `json:"ip,omitempty"`
	TalkGroup bool              `json:"talkgroup,omitempty"`
	Action    configEventAction `json:"action,omitempty"`
}

type configEventAction struct {
	Run  string   `json:"run,omitempty"`
	Args []string `json:"args,omitempty"`
}

// Global instance of structure
var config configData

// Load configuration
func configure(fileName string) error {

	// Load from json file
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	// Unmarshal data into config struct
	err = json.Unmarshal(file, &config)
	if err != nil {
		tmp := fmt.Sprintf("error parsing JSON: %s", err.Error())
		return errors.New(tmp)
	}

	// Iterate through clients and add
	for _, c := range config.Clients {
		clientAdd(c.ID, c.Password)
		log.Printf("Loaded client %d", c.ID)
	}

	// Iterate through events and log
	for _, e := range config.Events {
		if e.Enabled {
			log.Printf("Loaded event %s src %d dst %d client %d talkgroup %v ip %s action %s",
				e.Name, e.SRC, e.DST, e.Client, e.TalkGroup, e.IP, e.Action)
		} else {
			if config.Debug {
				log.Printf("Event %s is not enabled, ignoring", e.Name)
			}
		}
	}
	return nil
}
