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
	stream     uint32
	seq        uint32
	src        uint32
	dst        uint32
	repeater   uint32
	bitmap     uint32
	slot       uint32
	group      bool
	private    bool
	frameType  uint32
	frameVoice bool
	frameData  bool
	dataType   uint32
	data       bytes.Bytes
	addr       net.Addr
	ip         string
	original   bytes.Bytes
}

func DMRDParse(b bytes.Bytes, addr net.Addr) DMRData {
	var d DMRData
	d.addr = addr
	d.ip = strings.Split(safeAddrString(addr), ":")[0]
	d.seq = b.GetUint32(4, 1)
	d.src = b.GetUint32(5, 3)
	d.dst = b.GetUint32(8, 3)
	d.repeater = b.GetUint32(11, 4)
	d.bitmap = b.GetUint32(15, 1)
	d.stream = b.GetUint32(16, 4)
	d.data = b.Get(20, 0)
	d.original = b.Copy()

	// Get slot
	if d.bitmap&0x80 == 0x80 {
		d.slot = 2
	} else {
		d.slot = 1
	}

	// Call type - Private Call (PC) vs Talkgroup (TG)
	if d.bitmap&0x40 == 0x40 {
		d.group = false
		d.private = true
	} else {
		d.group = true
		d.private = false
	}

	// Frame Type
	d.frameType = (d.bitmap & 0x30) >> 4
	if d.frameType&0b10 == 0b10 {
		d.frameData = true
		d.frameVoice = false
	} else {
		d.frameData = false
		d.frameVoice = true
	}

	// Data Type
	d.dataType = d.bitmap & 0x0f

	// Print content of data frames for debugging
	//if d.frameType == 2 {
	//log.Printf("FrameType %02b, DT/VS %04b", d.frameType, d.dataType)
	//dump(d.data)
	//}

	return d
}

func DMRDSummary(c string, d DMRData) {
	log.Printf("%s DMRD seq=%d src=%d dst=%d slot=%v TG=%v PC=%v FT=%02b FV=%v FD=%v DT=%04b stream=%d repeater=%d @ %s\n",
		c, d.seq, d.src, d.dst, d.slot, d.group, d.private, d.frameType,
		d.frameVoice, d.frameData, d.dataType, d.stream, d.repeater, safeAddrString(d.addr))
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
