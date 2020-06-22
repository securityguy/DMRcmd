/*
	Copyright (c) 2020 by Eric Jacksch VE3XEJ

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.

    You should have received a copy of the GNU General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package main

import (
	"log"

	"dmrcmd/bytes"
)

// Send NAC to client (MSTNAK + id)
func sendNAK(dg datagram) {
	log.Printf("Sending RPTNAK to %d @ %s\n", dg.client.Uint32(), dg.addr.String())
	reply := bytes.New()
	reply.AppendString("MSTNAK")
	reply.Append(dg.client)
	sendUDP(dg, reply)
}

// Send ACK to client (RPTACK + id)
func sendACK(dg datagram) {
	log.Printf("Sending RPTACK to %d @ %s\n", dg.client.Uint32(), dg.addr.String())
	reply := bytes.New()
	reply.AppendString("RPTACK")
	reply.Append(dg.client)
	sendUDP(dg, reply)
}

// Send ping reply (pong) to client
func sendPONG(dg datagram) {
	log.Printf("Sending pong to %d @ %s\n", dg.client.Uint32(), dg.addr.String())
	reply := bytes.New()
	reply.AppendString("MSTPONG")
	reply.Append(dg.client)
	sendUDP(dg, reply)
}

// Send UDP datagram
func sendUDP(dg datagram, buf bytes.Bytes) {
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
