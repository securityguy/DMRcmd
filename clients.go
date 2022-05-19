/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/
package main

import (
	"crypto/sha256"
	"dmrcmd/bytes"
	"errors"
)

// Structure to hold information about clients
type client struct {
	id            uint32
	authenticated bool
	password      string
	salt          bytes.Bytes
	ip            string
}

// Create map
var clientList = make(map[uint32]client)

// Add client
// This can also be used to change the client's password
// and force a re-authentication
func clientAdd(id uint32, password string) {
	var c client
	c.id = id
	c.authenticated = false
	c.password = password
	c.salt = bytes.Bytes{}
	c.ip = ""
	clientList[id] = c
}

// Returns true if ID exists in map
func clientExist(id uint32) bool {
	if _, ok := clientList[id]; ok {
		return true
	}
	return false
}

// Store salt sent to client
func clientSalt(id uint32, salt bytes.Bytes) error {
	if clientExist(id) {
		c := clientList[id]
		c.salt = salt
		clientList[id] = c
		return nil
	}
	return errors.New("client ID does not exist")
}

// Authenticate client
func clientAuthenticate(id uint32, auth bytes.Bytes, ip string) bool {
	var result = false

	if !clientExist(id) {
		return false
	}

	// Get client record since we will be changing it
	c := clientList[id]

	// Calculate expected authentication
	hash := sha256.New()
	hash.Write(c.salt)
	hash.Write([]byte(c.password))
	expected := hash.Sum(nil)

	// Compare
	if auth.Equal(expected) {
		// success
		result = true
		c.ip = ip
	}

	// Salt can only be used once
	c.salt = bytes.Bytes{}

	// Save result for later queries
	c.authenticated = result
	clientList[id] = c
	return result
}

// Check if client is authenticated
func clientCheck(id uint32, ip string) bool {
	if !clientExist(id) {
		return false
	}

	if clientList[id].authenticated == true && clientList[id].ip == ip {
		return true
	}
	return false
}
