/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"fmt"
	"log"

	"dmrcmd/bytes"
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
		hex = hex + fmt.Sprintf("%02x ", b)

		// Combine at the right time
		if count%16 == 15 {
			log.Printf("%4d: %-47s | %s", address, hex, text)
			address = count + 1
			hex = ""
			text = ""
		}
	}
	log.Printf("%4d: %-48s | %s", address, hex, text)
}
