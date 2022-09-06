/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"dmrcmd/hotspot"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"dmrcmd/ha"
)

// Structure to hold configuration information
type configData struct {
	Debug    bool              `json:"debug"`
	Minimum  uint32            `json:"minimum,omitempty"`
	HA       ha.Config         `json:"ha,omitempty"`
	Hotspots []hotspot.Hotspot `json:"hotspots"`
	Events   []configEvent     `json:"events"`
}

type configEvent struct {
	Enabled   bool              `json:"enabled"`
	Name      string            `json:"name,omitempty"`
	SRC       uint32            `json:"src,omitempty"`
	DST       uint32            `json:"dst,omitempty"`
	Client    uint32            `json:"hotspot,omitempty"`
	IP        string            `json:"ip,omitempty"`
	TalkGroup bool              `json:"talkgroup,omitempty"`
	Action    configEventAction `json:"action,omitempty"`
}

type configEventAction struct {
	Run      string   `json:"run,omitempty"`
	Args     []string `json:"args,omitempty"`
	HAScript string   `json:"ha_script,omitempty"`
	HAScene  string   `json:"ha_scene,omitempty"`
}

// Global instance of structure
var config configData

// Load configuration
func configure(fileName string) error {

	// Load from json file
	file, err := os.ReadFile(fileName)
	if err != nil {
		return err
	}

	// Unmarshal data into config struct
	err = json.Unmarshal(file, &config)
	if err != nil {
		tmp := fmt.Sprintf("error parsing JSON: %s", err.Error())
		return errors.New(tmp)
	}

	// Iterate through hotspots and add
	for _, h := range config.Hotspots {
		if h.Enabled {
			hotspot.Add(h)
			log.Printf("Added hotspot %s [%d]", h.Name, h.ID)
			if len(h.Drop) > 0 {
				log.Printf("Hotspot %s [%d] configured to drop: %v", h.Name, h.ID, h.Drop)
			}
		} else {
			log.Printf("Ignoring disabled hotspot %s [%d]", h.Name, h.ID)
		}
	}

	// Iterate through events and log
	for _, e := range config.Events {
		if e.Enabled {
			log.Printf("Loaded event %s src %d dst %d hotspot %d talkgroup %v ip %s action %s",
				e.Name, e.SRC, e.DST, e.Client, e.TalkGroup, e.IP, actionToString(e.Action))
		} else {
			if config.Debug {
				log.Printf("Event %s is not enabled, ignoring", e.Name)
			}
		}
	}
	return nil
}
