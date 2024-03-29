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

// dmrProxy sits between a repeater and a DMR server
func dmrProxy(h hotspot.Hotspot) {

	// Destination for packet forwarding
	var dest net.Addr = nil
	var client net.Addr = nil
	var server net.Addr = nil

	// Listen for incoming udp packets
	pc, err := net.ListenPacket("udp", h.Listen)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Proxy listening on %s for repeater %s [%d] and forwarding to %s",
		h.Listen, h.Name, h.ID, h.Server)

	//noinspection GoUnhandledErrorResult
	defer pc.Close()

	// Set last purge to current time
	lastPurge := time.Now().Unix()

	// Get server address
	server, err = net.ResolveUDPAddr("udp", h.Server)
	if err != nil {
		log.Printf("Error parsing server address, terminating proxy for repeater %s [%d]: %s",
			h.Name, h.ID, err.Error())
		return
	}

	// Loop and receive UDP datagrams
	for {

		// ReadFrom will respect the length of buf, so we don't need to worry about buffer
		// overflows. If the packet contains more data than len(buf) it will be truncated.
		buf := make([]byte, 1024)
		n, addr, err := pc.ReadFrom(buf)
		if err != nil {
			continue
		}

		// We don't know what port the repeater will use, so we need to keep track of it
		if isLocal(addr) {
			// Make sure this isn't from the server just in case
			if safeAddrString(addr) != h.Server {
				// Is this a change?
				if safeAddrString(client) != safeAddrString(addr) {
					client = addr
					log.Printf("Client address set to %s for repeater %s [%d]",
						safeAddrString(client), h.Name, h.ID)
				}
			}
		}

		// If we haven't obtained a client address yet, don't proceed
		if client == nil {
			continue
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
		if safeAddrString(addr) == safeAddrString(client) {
			dg.local = true
		}

		if config.Debug {
			log.Printf("Received %d bytes from %s", n, safeAddrString(dg.addr))
			dump(dg.data)
		}

		// Process the datagram
		dispatchMsg(&dg)

		if dg.drop {
			if config.Debug {
				log.Printf("Proxy dropping %d bytes from %s", n, safeAddrString(dg.addr))
			}
			continue
		}

		// Determine destination
		if safeAddrString(addr) == h.Server {
			dest = client
		} else {
			dest = server
		}

		// Send the datagram
		n, err = pc.WriteTo(dg.data, dest)
		if err != nil {
			log.Printf("Send error on proxy for repeater %s [%d]: %s",
				h.Name, h.ID, err.Error())
			continue
		}

		if config.Debug {
			log.Printf("Sent %d bytes to %s", n, safeAddrString(dest))
		}

		// Get current time
		now := time.Now().Unix()

		// Purge old steams from the map every minute
		if now-lastPurge > 60 {
			eventPurge()
		}
	}
}
