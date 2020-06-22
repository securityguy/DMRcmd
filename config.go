/*
	Copyright (c) 2020 by Eric Jacksch VE3XEJ

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
