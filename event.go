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
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"
)

// Number of seconds after which to discard stream
const maxStreamAge int64 = 60

// Map of sequence numbers to identify new call
var stream = make(map[uint32]streamData)

type streamData struct {
	id        uint32
	last      int64
	triggered bool
	count     uint32
}

type dmrData struct {
	stream    uint32
	seq       uint32
	src       uint32
	dst       uint32
	client    uint32
	slot      uint32
	ip        string
	group     bool
	private   bool
	dataSync  bool
	voiceSync bool
}

// Filter to identify new transmissions and fire events
func eventFilter(d dmrData) {
	var s streamData
	now := time.Now().Unix()

	// Check if we have seen this stream before
	if _, ok := stream[d.stream]; ok {

		// Yes, we have seen it before
		s = stream[d.stream]
		s.count++

		// Have we already triggered on it?
		if stream[d.stream].triggered {
			s.last = now
			stream[d.stream] = s
			return
		}
	} else {
		// First time seeing this stream
		s.id = d.stream
		s.count = 1
		s.last = now
		s.triggered = false
	}

	if config.Debug {
		log.Printf("Minimum %d, count %d", config.Minimum, s.count)
	}

	// Is this stream long enough to trigger an event?
	// This is controlled by "seconds" in the global config
	if config.Minimum > 0 {
		if s.count < config.Minimum {
			// Not long enough yet
			stream[d.stream] = s
			return
		}
	}

	if config.Debug {
		log.Println("Reached minimum DMR data frame count, processing events")
	}

	// Stream will now trigger applicable events
	s.triggered = true

	// Iterate through events and look matches
	for _, c := range config.Events {

		// Ignore events that are disabled
		if c.Enabled == false {
			continue
		}

		// Only trigger on TGs if config talkgroup is true
		if c.TalkGroup == false && d.group == true {
			continue
		}

		// Match src if specified
		if c.SRC > 0 {
			if c.SRC != d.src {
				continue
			}
		}

		// Match dst if specified
		if c.DST > 0 {
			if c.DST != d.dst {
				continue
			}
		}

		// Match src if specified
		if c.Client > 0 {
			if c.Client != d.client {
				continue
			}
		}

		// Match IP if specified
		if c.IP != "" {
			if c.IP != d.ip {
				continue
			}
		}

		// All criteria met, perform action
		log.Printf("Triggered event %s from=%d to=%d client=%d private=%v group=%v ip=%s action: %s",
			c.Name, d.src, d.dst, d.client, d.private, d.group, d.ip, c.Action)
		eventAction(d, c.Action)
	}

	// Store updated data
	stream[d.stream] = s
}

// Perform action
func eventAction(d dmrData, action configEventAction) {

	// Build argument string starting with command
	var args []string
	args = append(args, action.Run)
	text := action.Run

	// Iterate through arguments
	for _, a := range action.Args {

		// Make substitutions
		arg := ""
		switch strings.ToLower(a) {

		case "$src":
			arg = fmt.Sprint(d.src)
		case "$dst":
			arg = fmt.Sprint(d.dst)
		case "$client":
			arg = fmt.Sprint(d.client)
		case "$ip":
			arg = d.ip
		default:
			arg = a
		}

		args = append(args, arg)
		text = text + " " + arg
	}

	// Build command and execute
	cmd := &exec.Cmd{
		Path: action.Run,
		Args: args,
	}
	err := cmd.Start()
	if err != nil {
		log.Printf("Error executing \"%s\": %s", text, err.Error())
		return
	}
	log.Printf("Executed \"%s\"", text)
}

// Purge old steams from the map
// Otherwise we will end storing every stream we see
func eventPurge() {
	purge := time.Now().Unix() - maxStreamAge

	for index, s := range stream {
		if s.last < purge {
			//if config.Debug {
			if true {
				log.Printf("Purging stream %d", index)
			}
			delete(stream, index)
		}
	}
}
