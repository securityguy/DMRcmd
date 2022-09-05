/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"
	"net"
	"time"

	"dmrcmd/bytes"
	"dmrcmd/hotspot"
)

// TODO -- turn this into a proxy instead of a server
func runProxy(h hotspot.Hotspot) {

	// Listen for incoming udp packets
	pc, err := net.ListenPacket("udp", h.Listen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for packets on %s for hotspot %s [%d]", h.Listen, h.Name, h.ID)

	//noinspection GoUnhandledErrorResult
	defer pc.Close()

	// Set last purge to current time
	lastPurge := time.Now().Unix()

	// Loop and receive UDP datagrams
	for {
		// ReadFrom will respect the length of buf, so we don't need to worry about buffer
		// overflows. If the packet contains more data than len(buf) it will be truncated.
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		// Create and populate structure
		dg := datagram{
			pc:     pc,
			addr:   addr,
			data:   buf[:n],
			client: bytes.New(),
		}

		if config.Debug {
			log.Printf("Received %d bytes from %s", n, dg.addr.String())
			dump(dg.data)
		}

		// Handle the datagram
		dispatch(dg)

		// Get current time
		now := time.Now().Unix()

		// Purge old steams from the map every minute
		if now-lastPurge > 60 {
			eventPurge()
		}
	}
}
