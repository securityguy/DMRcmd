/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/
package main

import (
	"log"
	"math/rand"
	"strings"

	"dmrcmd/bytes"
)

// RPTL - Login command is RPTL followed by 4 byte ID
func RPTL(dg datagram) {
	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()

	log.Printf("Connection request from %d @ %s\n", id, dg.addr.String())

	if clientExist(id) == false {
		log.Printf("Unknown client %d", id)
		sendNAK(dg)
		return
	}

	// Create 4 random bytes for salt (32-bit integer)
	salt := bytes.New()
	salt.AppendUint32(rand.Uint32())

	// Store salt
	err := clientSalt(id, salt)
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
func RPTK(dg datagram) {
	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()
	authBytes := dg.data.Get(8, 32)
	log.Printf("Authentication request from %d @ %s\n", id, dg.addr.String())

	if clientExist(id) == false {
		log.Printf("Unknown client %d", id)
		sendNAK(dg)
		return
	}

	if clientAuthenticate(id, authBytes, dg.addr.String()) {
		log.Printf("Authenticated %d @ %s\n", id, dg.addr.String())
		sendACK(dg)
	} else {
		log.Printf("Authentication failed for %d @ %s\n", id, dg.addr.String())
		sendNAK(dg)
	}
}

// RPTC - Configuration message
func RPTC(dg datagram) {
	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()
	log.Printf("Configuration from %d @ %s\n", id, dg.addr.String())

	if clientCheck(id, dg.addr.String()) {
		// Send ack
		sendACK(dg)
	} else {
		// Send nak
		sendNAK(dg)
	}
}

// RPTPING - Ping
func RPTPING(dg datagram) {
	dg.client = dg.data.Get(7, 4)
	id := dg.client.Uint32()
	log.Printf("Ping from %d @ %s\n", id, dg.addr.String())

	// If client is authenticated, reply
	if clientCheck(id, dg.addr.String()) {
		sendPONG(dg)
	} else {
		log.Printf("Client %d @ %s is not authenticated\n", id, dg.addr.String())
		sendNAK(dg)
	}
}

// DMRD - DMR Data
func DMRD(dg datagram) {
	var d dmrData
	d.seq = dg.data.GetUint32(4, 1)
	d.src = dg.data.GetUint32(5, 3)
	d.dst = dg.data.GetUint32(8, 3)
	dg.client = dg.data.Get(11, 4)
	bitmap := dg.data.GetUint32(15, 1)
	d.stream = dg.data.GetUint32(16, 4)

	d.client = dg.client.Uint32()
	d.ip = strings.Split(dg.addr.String(), ":")[0]

	// Get slot
	if bitmap&0x80 == 0x80 {
		d.slot = 2
	} else {
		d.slot = 1
	}

	// Private Call (PC) vs Talkgroup (TG)
	if bitmap&0x40 == 0x40 {
		d.group = false
		d.private = true
	} else {
		d.group = true
		d.private = false
	}

	// Data Sync
	if bitmap&0x20 == 0x20 {
		d.dataSync = true
	} else {
		d.dataSync = false
	}

	// Voice Sync
	if bitmap&0x10 == 0x10 {
		d.voiceSync = true
	} else {
		d.voiceSync = false
	}

	// Is this from an authenticated client?
	if !clientCheck(d.client, dg.addr.String()) {
		log.Printf("Igoring DMRD from unauthenticated %d @ %s\n", d.client, dg.addr.String())
		return
	}

	if config.Debug {
		log.Printf("DMRD seq=%d src=%d dst=%d slot=%v TG=%v PC=%v stream=%d from=%d @ %s\n",
			d.seq, d.src, d.dst, d.slot, d.group, d.private, d.stream, d.client, dg.addr.String())
	}

	// Send for event detection
	eventFilter(d)
}

// DMRA - DMR Talker Alias
func DMRA(dg datagram) {
	dg.client = dg.data.Get(4, 4)
	id := dg.client.Uint32()
	radio := dg.data.GetUint32(8, 3)
	alias := dg.data.GetString(13, 6)

	// Is this from an authenticated client?
	if !clientCheck(id, dg.addr.String()) {
		log.Printf("Igoring DMRA from unauthenticated %d @ %s\n", id, dg.addr.String())
		return
	}

	log.Printf("DMRA radio %d alias %s from %d @ %s\n", radio, alias, id, dg.addr.String())
}
