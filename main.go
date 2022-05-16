/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ

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
	"net"
	"os"
	"time"

	"dmrcmd/bytes"
)

const ProductName = "dmrcmd"
const ProductVersion = "0.0.2"

// Structure to hold UDP message metadata and contents
// This avoids having to pass multiple variables
type datagram struct {
	pc     net.PacketConn
	addr   net.Addr
	data   bytes.Bytes
	client bytes.Bytes
}

func main() {

	// Say hello
	log.Printf("Started %s version %s", ProductName, ProductVersion)

	// Load configuration - optionally using command line argument
	configFile := "dmrcmd.conf"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	err := configure(configFile)
	if err != nil {
		log.Fatal("CONFIG ERROR: ", err)
	}
	log.Printf("Loaded configuration from %s", configFile)

	// Listen for incoming udp packets
	pc, err := net.ListenPacket("udp", config.Listen)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Listening for packets on %s", config.Listen)

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

		// Purge old steams from the map every minute
		now := time.Now().Unix()
		if now-lastPurge > 60 {
			eventPurge()
		}
	}
}

// Send packet to appropriate function in messages.go
func dispatch(dg datagram) {
	if dg.data.MatchStartString("RPTL") {
		RPTL(dg)
	} else if dg.data.MatchStartString("RPTK") {
		RPTK(dg)
	} else if dg.data.MatchStartString("RPTC") {
		RPTC(dg)
	} else if dg.data.MatchStartString("RPTPING") {
		RPTPING(dg)
	} else if dg.data.MatchStartString("DMRD") {
		DMRD(dg)
	} else if dg.data.MatchStartString("DMRA") {
		DMRA(dg)
	} else {
		log.Printf("Unknown packet type from %s", dg.addr.String())
	}
}
