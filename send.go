/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"

	"dmrcmd/bytes"
)

// Send NAC to hotspot (MSTNAK + id)
func sendNAK(dg *datagram) {
	log.Printf("Sending RPTNAK to %d @ %s\n", dg.hotspot.Uint32(), dg.addr.String())
	reply := bytes.New()
	reply.AppendString("MSTNAK")
	reply.Append(dg.hotspot)
	sendUDP(dg, reply)
}

// Send ACK to hotspot (RPTACK + id)
func sendACK(dg *datagram) {
	log.Printf("Sending RPTACK to %d @ %s\n", dg.hotspot.Uint32(), dg.addr.String())
	reply := bytes.New()
	reply.AppendString("RPTACK")
	reply.Append(dg.hotspot)
	sendUDP(dg, reply)
}

// Send ping reply (pong) to hotspot
func sendPONG(dg *datagram) {
	log.Printf("Pong to %d @ %s\n", dg.hotspot.Uint32(), dg.addr.String())
	reply := bytes.New()
	reply.AppendString("MSTPONG")
	reply.Append(dg.hotspot)
	sendUDP(dg, reply)
}

// Send UDP datagram
func sendUDP(dg *datagram, buf bytes.Bytes) {
	n, err := dg.pc.WriteTo(buf, dg.addr)
	if err != nil {
		log.Printf("error sending UDP datagram to %s\n", dg.addr.String())
		return
	}

	if config.Debug {
		log.Printf("Sent %d bytes to %s", n, dg.addr.String())
		dump(buf)
	}
}
