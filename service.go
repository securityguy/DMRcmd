/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"dmrcmd/bytes"
	"dmrcmd/hotspot"
	"log"
	"net"
)

// Structure to hold UDP message metadata and contents
// This avoids having to pass multiple variables
type datagram struct {
	pc     net.PacketConn
	addr   net.Addr
	data   bytes.Bytes
	client bytes.Bytes
}

func startService(id uint32) {

	// Check that hotspot entry exists
	if !hotspot.Exists(id) {
		log.Printf("Unable to start server, %d does not exist", id)
		return
	}

	// Get hotspot configuration
	h := hotspot.Get(id)

	// Start server or proxy
	if h.Proxy {
		runProxy(h)
	} else {
		runServer(h)
	}
}
