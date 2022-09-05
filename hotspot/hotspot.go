/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/
package hotspot

import (
	"crypto/sha256"
	"errors"
	"strings"

	"dmrcmd/bytes"
)

// Hotspot information
type Hotspot struct {
	Enabled       bool   `json:"enabled"`
	Name          string `json:"name"`
	ID            uint32 `json:"id"`
	Password      string `json:"password"`
	Mode          string `json:"mode"`
	Listen        string `json:"listen"`
	Server        string `json:"server"`
	authenticated bool
	proxy         bool
	salt          bytes.Bytes
	addr          string
}

// Create map
var hotspots = make(map[uint32]Hotspot)

// Add hotspot
// This can also be used to change the client's password
// and force a re-authentication
func Add(h Hotspot) {
	if strings.ToLower(h.Mode) == "proxy" {
		h.proxy = true
	} else {
		h.proxy = false
	}

	h.authenticated = false
	h.salt = bytes.Bytes{}
	h.addr = ""
	hotspots[h.ID] = h
}

// Exists returns true if ID exists in map
func Exists(id uint32) bool {
	if _, ok := hotspots[id]; ok {
		return true
	}
	return false
}

// Get configuration information for hotspot
func Get(id uint32) Hotspot {
	if _, ok := hotspots[id]; ok {
		return hotspots[id]
	}
	return Hotspot{}
}

// GetList configuration information for hotspot
func GetList() []uint32 {
	var list []uint32

	// Iterate over hotspots and build list
	for i, _ := range hotspots {
		list = append(list, i)
	}

	return list
}

// Salt stores the salt sent to hotspot in server mode
func Salt(id uint32, salt bytes.Bytes) error {
	if Exists(id) {
		h := hotspots[id]
		h.salt = salt
		hotspots[id] = h
		return nil
	}
	return errors.New("hotspot ID does not exist")
}

// Authenticate hotspot
func Authenticate(id uint32, auth bytes.Bytes, ip string) bool {
	var result = false

	if !Exists(id) {
		return false
	}

	// Get client record since we will be changing it
	h := hotspots[id]

	// Calculate expected authentication
	hash := sha256.New()
	hash.Write(h.salt)
	hash.Write([]byte(h.Password))
	expected := hash.Sum(nil)

	// Compare
	if auth.Equal(expected) {
		// success
		result = true
		h.addr = ip
	}

	// Salt can only be used once
	h.salt = bytes.Bytes{}

	// Save result for later queries
	h.authenticated = result
	hotspots[id] = h
	return result
}

// Check if client is authenticated. Note that addr contains the IP address and port.
func Check(id uint32, addr string) bool {
	if !Exists(id) {
		return false
	}

	if hotspots[id].authenticated == true && hotspots[id].addr == addr {
		return true
	}
	return false
}
