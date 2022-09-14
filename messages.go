/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"
	"math/rand"

	"dmrcmd/bytes"
	"dmrcmd/hotspot"
)

// Send packet to appropriate message handler
func dispatchMsg(dg *datagram) {
	if dg.data.MatchStartString("RPTL") {
		RPTL(dg)
	} else if dg.data.MatchStartString("RPTK") {
		RPTK(dg)
	} else if dg.data.MatchStartString("RPTC") {
		RPTC(dg)
	} else if dg.data.MatchStartString("RPTPING") {
		RPTPING(dg)
	} else if dg.data.MatchStartString("MSTPONG") {
		MSTPONG(dg)
	} else if dg.data.MatchStartString("DMRD") {
		DMRD(dg)
	} else if dg.data.MatchStartString("DMRA") {
		DMRA(dg)
	} else {
		if !dg.proxy {
			log.Printf("Unknown packet type from %s", safeAddrString(dg.addr))
		}
	}
}

// RPTL - Login command is RPTL followed by 4 byte ID
func RPTL(dg *datagram) {

	// If proxying, take no action
	if dg.proxy {
		return
	}

	dg.hotspot = dg.data.Get(4, 4)
	id := dg.hotspot.Uint32()

	log.Printf("Connection request from %d @ %s\n", id, safeAddrString(dg.addr))

	if hotspot.Exists(id) == false {
		log.Printf("Unknown repeater %d", id)
		sendNAK(dg)
		return
	}

	// Create 4 random bytes for salt (32-bit integer)
	salt := bytes.New()
	salt.AppendUint32(rand.Uint32())

	// Store salt
	err := hotspot.Salt(id, salt)
	if err != nil {
		log.Printf("error adding salt to repeater record")
		sendNAK(dg)
		return
	}

	// Reply with RPTACK + salt
	reply := bytes.New()
	reply.AppendString("RPTACK")
	reply.Append(salt)
	sendUDP(dg, reply)
}

// RPTK - Authentication
func RPTK(dg *datagram) {

	// If proxying, take no action
	if dg.proxy {
		return
	}

	dg.hotspot = dg.data.Get(4, 4)
	id := dg.hotspot.Uint32()
	authBytes := dg.data.Get(8, 32)
	log.Printf("Authentication request from %d @ %s\n", id, safeAddrString(dg.addr))

	if hotspot.Exists(id) == false {
		log.Printf("Unknown repeater %d", id)
		sendNAK(dg)
		return
	}

	if hotspot.Authenticate(id, authBytes, safeAddrString(dg.addr)) {
		log.Printf("Authenticated %d @ %s\n", id, safeAddrString(dg.addr))
		sendACK(dg)
	} else {
		log.Printf("Authentication failed for %d @ %s\n", id, safeAddrString(dg.addr))
		sendNAK(dg)
	}
}

// RPTC - Configuration message
func RPTC(dg *datagram) {

	// If proxying, take no action
	if dg.proxy {
		return
	}

	dg.hotspot = dg.data.Get(4, 4)
	id := dg.hotspot.Uint32()
	log.Printf("Configuration from %d @ %s\n", id, safeAddrString(dg.addr))

	if hotspot.CheckAuthenticated(id, safeAddrString(dg.addr)) {
		// Send ack
		sendACK(dg)
	} else {
		// Send nak
		sendNAK(dg)
	}
}

// RPTPING - Ping
func RPTPING(dg *datagram) {
	dg.hotspot = dg.data.Get(7, 4)
	id := dg.hotspot.Uint32()
	log.Printf("Ping from %d at %s\n", id, safeAddrString(dg.addr))

	// If proxying, take no further action
	if dg.proxy {
		return
	}

	// If repeater is authenticated, reply
	if hotspot.CheckAuthenticated(id, safeAddrString(dg.addr)) {
		sendPONG(dg)
	} else {
		log.Printf("Client %d @ %s is not authenticated\n", id, safeAddrString(dg.addr))
		sendNAK(dg)
	}
}

// MSTPONG - Pong
func MSTPONG(dg *datagram) {
	dg.hotspot = dg.data.Get(7, 4)
	id := dg.hotspot.Uint32()
	log.Printf("Pong to %d from %s\n", id, safeAddrString(dg.addr))
}

// DMRD - DMR Data
func DMRD(dg *datagram) {
	d := DMRDParse(dg.data, dg.addr)

	if dg.proxy {
		// Set drop flag if required
		dg.drop = hotspot.CheckDrop(d.repeater, d.src, d.dst)

		// CheckAuthenticated if datagram is from a local repeater
		if dg.local == false {
			if config.Debug {
				log.Printf("Not processing DMRD because source is not local %d @ %s\n", d.repeater, safeAddrString(dg.addr))
			}
			return
		}
	} else {
		// CheckAuthenticated if repeater has authenticated
		if !hotspot.CheckAuthenticated(d.repeater, safeAddrString(dg.addr)) {
			log.Printf("Ignoring DMRD from unauthenticated %d @ %s\n", d.repeater, safeAddrString(dg.addr))
			return
		}
	}

	if config.Debug {
		DMRDSummary("R", d)
	}

	// Send for event detection
	eventFilter(d)
}

// DMRA - DMR Talker Alias
func DMRA(dg *datagram) {
	dg.hotspot = dg.data.Get(4, 4)
	id := dg.hotspot.Uint32()
	radio := dg.data.GetUint32(8, 3)
	alias := dg.data.GetString(13, 6)

	// Is this from an authenticated repeater?
	if hotspot.CheckAuthenticated(id, safeAddrString(dg.addr)) || dg.proxy {
		log.Printf("DMRA radio %d alias %s from %d @ %s\n", radio, alias, id, safeAddrString(dg.addr))
	} else {
		log.Printf("Igoring DMRA from unauthenticated %d @ %s\n", id, safeAddrString(dg.addr))
	}
}
