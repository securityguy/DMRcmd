/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"dmrcmd/ha"
)

// Number of seconds after which to discard stream
const maxStreamAge int64 = 60

// Map of sequence numbers to identify new call
var stream = make(map[uint32]streamData)

type streamData struct {
	id         uint32
	last       int64
	voiceCount uint32
	dataCount  uint32
	considered []int
}

// Filter to identify new transmissions and trigger events
func eventFilter(d DMRData) {
	var s streamData
	now := time.Now().Unix()

	// Check if we have seen this stream before
	if _, ok := stream[d.stream]; ok {
		// Yes, we have seen it before
		s = stream[d.stream]

	} else {
		// First time seeing this stream
		s.id = d.stream
		s.voiceCount = 0
		s.dataCount = 0
	}

	// Update last time seen
	// this is used for garbage collection
	s.last = now

	// Increment data or voice frame counter
	if d.frameData {
		s.dataCount++
	} else {
		s.voiceCount++
	}

	if config.Debug {
		log.Printf("Voice frames: %d, Data frames: %d", s.voiceCount, s.dataCount)
	}

	// Iterate through events and look for events to trigger
	for i, c := range config.Events {

		// Ignore events that are disabled
		if c.Enabled == false {
			continue
		}

		// Check if we have previously reached the trigger threshold for this event
		done := false
		for _, t := range s.considered {
			if t == i {
				// The required number of frames have previously been
				// reached for this event and it has been checked below.
				// Abort to avoid repeating events
				done = true
				break
			}
		}

		if done {
			continue
		}

		// check for required minimum number of frames
		requiredFrames := false
		if c.RequiredVoice > 0 && s.voiceCount >= c.RequiredVoice {
			requiredFrames = true
		}

		if c.RequiredData > 0 && s.dataCount >= c.RequiredData {
			requiredFrames = true
		}

		if requiredFrames == false {
			continue
		} else {
			// Minimum number of frames reached, so no matter what happens
			// after this point, we don't want to trigger again
			if config.Debug {
				log.Printf("Reached minimum frame count for event \"%s\"", c.Name)
			}
			s.considered = append(s.considered, i)
		}

		// Continue checking other parameters

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
			if c.Client != d.repeater {
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
		log.Printf("Triggered event %s from=%d to=%d repeater=%d private=%v group=%v ip=%s action: %s",
			c.Name, d.src, d.dst, d.repeater, d.private, d.group, d.ip, actionToString(c.Action))
		go eventAction(d, c.Action)
	}

	// Store updated data
	stream[d.stream] = s
}

// Perform action
func eventAction(d DMRData, action configEventAction) {

	// Is there a command to execute?
	if action.Run != "" {
		eventExecute(d, action)
	}

	if action.HAScript != "" {
		h := ha.New(config.HA)
		err := h.Script(action.HAScript)
		if err != nil {
			log.Printf("Error triggering Home Assistant script \"%s\": %s", action.HAScript, err.Error())
		}
		log.Printf("Successfully considered Home Assistant script \"%s\"", action.HAScript)
	}

	if action.HAScene != "" {
		h := ha.New(config.HA)
		err := h.Scene(action.HAScene)
		if err != nil {
			log.Printf("Error triggering Home Assistant script \"%s\": %s", action.HAScene, err.Error())
		}
		log.Printf("Successfully considered Home Assistant script \"%s\"", action.HAScene)
	}
}

// Execute command
func eventExecute(d DMRData, action configEventAction) {

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
			arg = fmt.Sprint(d.repeater)
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

func actionToString(c configEventAction) string {
	a, err := json.Marshal(c)
	if err != nil {
		return ""
	}
	return string(a)
}
