/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"dmrcmd/bytes"
	"dmrcmd/hotspot"
	"log"
	"math/rand"
)

// RPTL - Login command is RPTL followed by 4 byte ID
func RPTL(dg *datagram) {

	// If proxying, take no action
	if dg.proxy {
		return
	}

	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()

	log.Printf("Connection request from %d @ %s\n", id, dg.addr.String())

	if hotspot.Exists(id) == false {
		log.Printf("Unknown client %d", id)
		sendNAK(dg)
		return
	}

	// Create 4 random bytes for salt (32-bit integer)
	salt := bytes.New()
	salt.AppendUint32(rand.Uint32())

	// Store salt
	err := hotspot.Salt(id, salt)
	if err != nil {
		log.Printf("error adding salt to client record")
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

	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()
	authBytes := dg.data.Get(8, 32)
	log.Printf("Authentication request from %d @ %s\n", id, dg.addr.String())

	if hotspot.Exists(id) == false {
		log.Printf("Unknown client %d", id)
		sendNAK(dg)
		return
	}

	if hotspot.Authenticate(id, authBytes, dg.addr.String()) {
		log.Printf("Authenticated %d @ %s\n", id, dg.addr.String())
		sendACK(dg)
	} else {
		log.Printf("Authentication failed for %d @ %s\n", id, dg.addr.String())
		sendNAK(dg)
	}
}

// RPTC - Configuration message
func RPTC(dg *datagram) {

	// If proxying, take no action
	if dg.proxy {
		return
	}

	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()
	log.Printf("Configuration from %d @ %s\n", id, dg.addr.String())

	if hotspot.CheckAuthenticated(id, dg.addr.String()) {
		// Send ack
		sendACK(dg)
	} else {
		// Send nak
		sendNAK(dg)
	}
}

// RPTPING - Ping
func RPTPING(dg *datagram) {
	dg.client = dg.data.Get(7, 4)
	id := dg.client.Uint32()
	log.Printf("Ping from %d at %s\n", id, dg.addr.String())

	// If proxying, take no further action
	if dg.proxy {
		return
	}

	// If client is authenticated, reply
	if hotspot.CheckAuthenticated(id, dg.addr.String()) {
		sendPONG(dg)
	} else {
		log.Printf("Client %d @ %s is not authenticated\n", id, dg.addr.String())
		sendNAK(dg)
	}
}

// MSTPONG - Pong
func MSTPONG(dg *datagram) {
	dg.client = dg.data.Get(7, 4)
	id := dg.client.Uint32()
	log.Printf("Pong to %d from %s\n", id, dg.addr.String())
}

// DMRD - DMR Data
func DMRD(dg *datagram) {
	d := DMRDParse(dg.data, dg.addr)

	if dg.proxy {
		// Set drop flag if required
		dg.drop = hotspot.CheckDrop(d.client, d.src, d.dst)

		// CheckAuthenticated if datagram is from a local hotspot
		if dg.local == false {
			if config.Debug {
				log.Printf("Not processing DMRD because source is not local %d @ %s\n", d.client, d.addr.String())
			}
			return
		}
	} else {
		// CheckAuthenticated if hotspot has authenticated
		if !hotspot.CheckAuthenticated(d.client, dg.addr.String()) {
			log.Printf("Ignoring DMRD from unauthenticated %d @ %s\n", d.client, d.addr.String())
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
	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()
	radio := dg.data.GetUint32(8, 3)
	alias := dg.data.GetString(13, 6)

	// Is this from an authenticated client?
	if hotspot.CheckAuthenticated(id, dg.addr.String()) || dg.proxy {
		log.Printf("DMRA radio %d alias %s from %d @ %s\n", radio, alias, id, dg.addr.String())
	} else {
		log.Printf("Igoring DMRA from unauthenticated %d @ %s\n", id, dg.addr.String())
	}
}
