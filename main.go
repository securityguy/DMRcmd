/*
		Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ

	    This program is free software: you can redistribute it and/or modify
	    it under the terms of the GNU General Public License as published by
	    the Free Software Foundation, either version 3 of the License, or
	    (at your option) any later version.

	    This program is distributed in the hope that it will be useful,
	    but WITHOUT ANY WARRANTY; without even the implied warranty of
	    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
	    GNU General Public License for more details.

	    You should have received a copy of the GNU General Public License
	    along with this program. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import (
	"dmrcmd/hotspot"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

const ProductName = "dmrcmd"
const ProductVersion = "0.0.12"

func main() {

	// Say hello
	log.Printf("Starting %s v%s on %s/%s",
		ProductName, ProductVersion, runtime.GOOS, runtime.GOARCH)

	// Load configuration - optionally using command line argument
	configFile := "dmrcmd.json"
	if len(os.Args) > 1 {
		configFile = os.Args[1]
	}

	err := configure(configFile)
	if err != nil {
		log.Fatal("CONFIG ERROR: ", err)
	}
	log.Printf("Loaded configuration from %s", configFile)

	// Setup signal catching
	signals := make(chan os.Signal, 1)

	// Catch signals
	signal.Notify(signals, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	// method invoked upon seeing signal
	go func() {
		for {
			s := <-signals
			log.Printf("Received signal: %v", s)
			AppCleanup()
		}
	}()

	// Start servers to handle each repeater
	list := hotspot.GetList()
	for _, id := range list {

		// Start server or proxy
		go startServer(id)
	}

	// Loop until terminated. Select uses less CPU than for{}
	select {}
}

func AppCleanup() {

	// Log exit
	log.Printf("Stopping %s v%s on %s/%s",
		ProductName, ProductVersion, runtime.GOOS, runtime.GOARCH)

	// Exit
	os.Exit(0)
}
