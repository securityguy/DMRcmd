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
	} else if dg.data.MatchStartString("DMRD") {
		DMRD(dg)
	} else if dg.data.MatchStartString("DMRA") {
		DMRA(dg)
	} else {
		log.Printf("Unknown packet type from %s", dg.addr.String())
	}
}
