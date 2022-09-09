/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"
	"net"

	"dmrcmd/bytes"
	"dmrcmd/hotspot"
)

// Structure to hold UDP message metadata and contents
// This avoids having to pass multiple variables
type datagram struct {
	pc      net.PacketConn // connection
	addr    net.Addr       // source address
	data    bytes.Bytes    // data
	hotspot bytes.Bytes    // hotspot ID
	proxy   bool           // proxy mode flag
	local   bool           // datagram is from local hotspot (proxy mode only)
	drop    bool           // drop packet flag (proxy mode only)
}

func startServer(id uint32) {

	// CheckAuthenticated that repeater entry exists
	if !hotspot.Exists(id) {
		log.Printf("Unable to start server, %d does not exist", id)
		return
	}

	// Get repeater configuration
	h := hotspot.Get(id)

	// Start server or proxy
	if h.Proxy {
		dmrProxy(h)
	} else {
		dmrServer(h)
	}
}
