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
	pc     net.PacketConn // connection
	addr   net.Addr       // source address
	data   bytes.Bytes    // data
	client bytes.Bytes    // TODO is this still required?
	proxy  bool           // proxy mode flag
	local  bool           // datagram is from local hotspot (proxy mode only)
	drop   bool           // drop packet flag (proxy mode only)
}

func startService(id uint32) {

	// CheckAuthenticated that hotspot entry exists
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
