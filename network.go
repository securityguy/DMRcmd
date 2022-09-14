/*
Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/

package main

import (
	"log"
	"net"
	"strings"
)

func isLocal(addr net.Addr) bool {
	return isInNetwork(addr, config.LocalNet)
}

func isInNetwork(addr net.Addr, cidr string) bool {

	// Parse CIDR
	_, localNet, err := net.ParseCIDR(cidr)
	if err != nil {
		log.Printf("Error parsing CIDR %s: %s", cidr, err.Error())
		return false
	}

	// split addr.string() into ip and port
	parts := strings.Split(addr.String(), ":")
	if len(parts) < 1 {
		log.Printf("Error splitting address into IP and port %s", addr.String())
		return false
	}

	ipAddress := net.ParseIP(parts[0])
	if ipAddress == nil {
		log.Printf("Error parsing IP %s", parts[0])
		return false
	}

	// Is IP within CIDR?
	if localNet.Contains(ipAddress) {
		return true
	}
	return false
}
