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
	Debug        bool              `json:"debug"`
	DefaultVoice uint32            `json:"default_voice,omitempty"`
	DefaultData  uint32            `json:"default_data,omitempty"`
	LocalNet     string            `json:"local_network,omitempty"`
	HA           ha.Config         `json:"ha,omitempty"`
	Hotspots     []hotspot.Hotspot `json:"hotspots"`
	Events       []configEvent     `json:"events"`
}

type configEvent struct {
	Enabled       bool              `json:"enabled"`
	Name          string            `json:"name,omitempty"`
	SRC           uint32            `json:"src,omitempty"`
	DST           uint32            `json:"dst,omitempty"`
	Client        uint32            `json:"repeater,omitempty"`
	IP            string            `json:"ip,omitempty"`
	TalkGroup     bool              `json:"talkgroup,omitempty"`
	RequiredData  uint32            `json:"required_data,omitempty"`
	RequiredVoice uint32            `json:"required_voice,omitempty"`
	Action        configEventAction `json:"action,omitempty"`
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
			log.Printf("Added repeater %s [%d]", h.Name, h.ID)
			if len(h.Drop) > 0 {
				log.Printf("Hotspot %s [%d] configured to drop: %v", h.Name, h.ID, h.Drop)
			}
		} else {
			log.Printf("Ignoring disabled repeater %s [%d]", h.Name, h.ID)
		}
	}

	// Iterate through events and log
	for i, e := range config.Events {
		if e.Enabled {
			// If required frames not set, use defaults
			if e.RequiredData == 0 {
				config.Events[i].RequiredData = config.DefaultData
			}

			if e.RequiredVoice == 0 {
				config.Events[i].RequiredVoice = config.DefaultVoice
			}

			log.Printf("Loaded event \"%s\" src %d dst %d repeater %d talkgroup %v ip %s action %s",
				e.Name, e.SRC, e.DST, e.Client, e.TalkGroup, e.IP, actionToString(e.Action))
		} else {
			if config.Debug {
				log.Printf("Event %s is not enabled, ignoring", e.Name)
			}
		}
	}

	// Sanity checks
	if config.LocalNet == "" {
		return errors.New("fatal configuration error, local_network is not defined")
	}

	return nil
}
