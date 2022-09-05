/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"
)

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
	} else if dg.data.MatchStartString("MSTPONG") {
		MSTPONG(dg)
	} else if dg.data.MatchStartString("DMRD") {
		DMRD(dg)
	} else if dg.data.MatchStartString("DMRA") {
		DMRA(dg)
	} else {
		if !dg.proxy {
			log.Printf("Unknown packet type from %s", dg.addr.String())
		}
	}
}
