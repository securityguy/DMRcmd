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
	"dmrcmd/bytes"
	"fmt"
	"log"
)

func dump(data bytes.Bytes) {
	hex := ""
	text := ""
	address := 0

	// Iterate over data
	for count, b := range data {

		// Text portion
		if b < 32 || b >= 127 {
			text = text + "."
		} else {
			text = text + string(b)
		}

		// Hex portion
		hex = hex + fmt.Sprintf("%02x", b)
		if count%4 == 3 {
			hex = hex + " "
		}

		//fmt.Printf("%d %d %d\n", totalCount, totalCount%4, totalCount%32)
		if count%32 == 31 {
			log.Printf("%4d: %-72s | %s", address, hex, text)
			address = count + 1
			hex = ""
			text = ""
		}
	}
	log.Printf("%4d: %-72s | %s", address, hex, text)
}
