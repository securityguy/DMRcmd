/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"
	"net"
	"strings"

	"dmrcmd/bytes"
)

type DMRData struct {
	stream    uint32
	seq       uint32
	src       uint32
	dst       uint32
	client    uint32
	slot      uint32
	group     bool
	private   bool
	dataSync  bool
	voiceSync bool
	data      bytes.Bytes
	addr      net.Addr
	ip        string
	original  bytes.Bytes
}

func DMRDParse(b bytes.Bytes, addr net.Addr) DMRData {
	var d DMRData
	d.addr = addr
	d.ip = strings.Split(addr.String(), ":")[0]
	d.seq = b.GetUint32(4, 1)
	d.src = b.GetUint32(5, 3)
	d.dst = b.GetUint32(8, 3)
	d.client = b.GetUint32(11, 4)
	bitmap := b.GetUint32(15, 1)
	d.stream = b.GetUint32(16, 4)
	d.data = b.Get(20, 0)
	d.original = b.Copy()

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

	return d
}

func DMRDSummary(c string, d DMRData) {
	log.Printf("%s DMRD seq=%d src=%d dst=%d slot=%v TG=%v PC=%v stream=%d client=%d @ %s\n",
		c, d.seq, d.src, d.dst, d.slot, d.group, d.private, d.stream, d.client, d.addr.String())
}

func DMRDDump(c string, d DMRData) {
	DMRDSummary(c, d)
	dump(d.original)
}

//goland:noinspection GoUnusedExportedFunction
func DMRDDumpRaw(c string, b bytes.Bytes, a net.Addr) {
	tmp := DMRDParse(b, a)
	DMRDDump(c, tmp)
}
