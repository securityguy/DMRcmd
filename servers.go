/*
	Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ
*/
package main

// Structure to hold information about clients
type server struct {
	name     string
	host     string
	id       uint32
	password string
	def      bool
}

// Create map
var serverList = make(map[string]server)

// Add server
func serverAdd(newServer configServer) {
	var s server
	s.name = newServer.Name
	s.host = newServer.Host
	s.id = newServer.ID
	s.password = newServer.Password
	s.def = newServer.Default
	serverList[s.host] = s
}
