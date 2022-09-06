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

// dmrProxy sits between a hotspot and a DMR server
func dmrProxy(h hotspot.Hotspot) {

	// Destination for packet forwarding
	var dest net.Addr = nil
	var client net.Addr = nil

	// Listen for incoming udp packets
	pc, err := net.ListenPacket("udp", h.Listen)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Proxy listening on %s for hotspot %s [%d] and forwarding to %s",
		h.Listen, h.Name, h.ID, h.Server)

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

		// We don't know what port the hotspot will use, so we need to keep track of it
		if client == nil {
			// Make sure this isn't from the server
			if addr.String() != h.Server {
				client = addr
				log.Printf("Client address set to %s for hotspot %s [%d]", client.String(), h.Name, h.ID)
			}
		}

		// Create and populate structure
		dg := datagram{
			pc:      pc,
			addr:    addr,
			data:    buf[:n],
			hotspot: bytes.New(),
			proxy:   true,
			local:   false,
			drop:    false,
		}

		// Set local flag if applicable
		if addr.String() == client.String() {
			dg.local = true
		}

		if config.Debug {
			log.Printf("Received %d bytes from %s", n, dg.addr.String())
			dump(dg.data)
		}

		// Process the datagram
		dispatchMsg(&dg)

		if dg.drop {
			if config.Debug {
				log.Printf("Proxy dropping %d bytes from %s", n, dg.addr.String())
			}
			continue
		}

		// Determine destination
		if addr.String() == h.Server {
			dest = client
		} else {
			dest, err = net.ResolveUDPAddr("udp", h.Server)
			if err != nil {
				log.Printf("Error parsing server address")
				break
			}
		}

		// Send the datagram
		n, err = pc.WriteTo(dg.data, dest)
		if err != nil {
			log.Printf("SEND ERROR!")
			continue
		}

		if config.Debug {
			log.Printf("Sent %d bytes to %s", n, dest.String())
		}

		// Get current time
		now := time.Now().Unix()

		// Purge old steams from the map every minute
		if now-lastPurge > 60 {
			eventPurge()
		}
	}
}
