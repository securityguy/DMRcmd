# DMRcmd

Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ.

This is early beta software. Please read this file in its entirety and use at your own risk.

This program acts as an amateur radio DMR network server or proxy using the Homebrew MMDVMHost protocol
and can be configured to execute commands based on the DMR traffic it receives from a hotspot.

In addition to executing command-line functions, DMRcmd also includes a Home Assistant integration.

This software was inspired by PiStar-Remote by Andy Taylor (MW0MWZ) and DMRGateway by Jonathan Naylor (G4KLX).
Collaboration and pull requests are welcome.

### Cautions

**Users must ensure that this program it not used in any way that places persons or
property at risk. It is not suitable for any application even remotely associated with life safety.
It should not be run on a computer containing sensitive or valuable information. By default, this program does not
do anything. The user is solely responsible for whatever they configure it to do.**

The Homebrew repeater protocol used by this application is not fully documented and as a result this implementation
from scratch may contain errors.

DMRcmd is intended for use on a private secure network. While it ignores traffic from unauthenticated sources,
DMR data packets are sent in unauthenticated UDP datagrams, making it easy to forge datagrams and trigger events.

The DMR protocol is not secure and this application can be configured to execute the system commands specified in the
configuration file. Anyone with a suitable receiver can monitor DMR traffic, view the source and destination IDs,
and program a radio to duplicate the traffic.

It is up to each user to determine the suitability of this software for their specific application and to accept
full responsibility for the outcome.

This software is intended for amateur radio and educational purposes only.

### Compiling

This program is written in Go (aka golang). To compile:

1) If you have not already done so, install Go from https://golang.org/ or using your package manager.
2) Download or clone the repo from https://github.com/securityguy/DMRcmd.
3) Change to the DMRcmd directory and type `go build`.

This software should compile and run on various flavours of Linux, Windows, and macOS.
If you encounter any errors please create an issue in GitHub.

Notes:

- Once compiled, the binary and configuration file can simply be copied to another computer.
- Go is only required to compile the software. It is not required to run the compiled program.
- Windows users should consider installing Git from https://git-scm.com/.

### Configuration

By default, configuration information is read from dmrcmd.json in the working directory. The full path to the configuration file can optionally be passed as the first (and only) command line argument.

Note that most configuration objects have an "enabled" option. Unless it is set to true, the object will be ignored.

You must update "local_network" in CIDR format (for example 192.168.0.0/24). Only hotspots within the specified network will be allowed to connect. In addition to the security benefits, this allows the proxy module to reliably track the hotspot address and port if it changes. The server module will also ignore UDP datagrams that do not originate from the configured local network.

Additional configuration documentation will be added at a later date. In the interim, please refer to dmrcmd.example.json.

For security reasons, this program neither searches the execution path nor uses a shell to execute commands. The full path to a program or script to be executed must be provided in the configuration file.

For similar reasons, command line arguments to be passed to the program or script must be specified as a list. Note that information from the DMR transmission (source ID, destination ID, client (repeater) ID, and IP address) can be substituted for command line arguments using $src, $dst, $ip, and $client respectively.

Within each event definition, **all specified conditions** must be met for the event to execute.

By design, a single DMR transmission may result in multiple events executing.

In order to avoid accidentally triggering events, a minimum number of DMR voice or data frames can be required.
The default number of required frames is specified in default_voice and default_data respectively. This can be
overridden on a per-event basis using required_voice and required_data.

The default value of requiring 10 voice frames requires the transmission to last for approximately one second
before triggering an event, which should ignore a quick bump of the PTT. By default, data frames are ignored (set to 0). Users should test the defaults and adjust them as required.

DMR emergency alerts consist of a single data frame. Therefore, to use a DMR emergency transmission to trigger an event, you will need to set required_data to 1 for the event you wish to trigger. Also note that each time the radio transmits the emergency alert it will be considered a separate event.

### Pi-Star and DRMGateway Users

Pi-Star contains a "DMR Gateway" that is capable of routing DMR traffic to multiple servers.
DMRGateway by Jonathan Naylor (G4KLK) is also available as a stand-alone application from
https://github.com/g4klx/DMRGateway.

The most straightforward approach to using DMRcmd is to use DMRGateway to route a range of DMR IDs to
DMRcmd, while allowing the remainder of your DMR traffic to flow to Brandmeister, DMR+, XLX, and other services as
usual.

For example, in the DMRGateway configuration file (Pi-Star users can access it at Configuration -> Expert -> Full Edit:
DMR GW):

    [DMR Network 2]
     Enabled=1
     Address=<IP where DMRcmd can be reached>
     Port=<PORT specified in the DMRcmd configuration file>
     PCRewrite=2,8999900,2,8999900,100
     TGRewrite=2,8999900,2,8999900,100
     Password=<PASSWORD>
     Debug=0
     Id=<DMR ID>
     Name=DMRcmd

This example will send any private or group calls in the 8999900 - 8999999 range to DMRcmd.
Using a Private Call is recommended. DMRcmd ignores group calls unless "talkgroup" is
set to true in the event configuration, in which case both private and group calls will be
considered for the event.

The DMR ID and password in the DMR Gateway configuration must match a hotspot configuration in DMRcmd
or authentication will fail and data from your hotspot will be ignored.

**Please be cognizant of the fact that any given DMR ID may be assigned to an individual and use appropriate
care to ensure that you do not inadvertently route traffic intended to trigger DMRcmd events to a DMR network.**

### openSPOT Users

openSPOT does not include multi-DMR server routing capability. If you wish to use multiple DMR networks, your best bet
is to connect openSPOT to DMRGateway (link above) and then use the same instructions as provided above to connect
DMRGateway to DMRcmd.

DMRcmd has a proxy option, in which case instead of acting as a DMR network server, it will proxy data
to and from the server specified in the configuration file. This allows DMRcmd to sit between OpenSPOT
and one DMR Network, and monitor the traffic for DMR calls of interest.

To configure openSPOT to send traffic to DMRcmd:

1) Log into openSPOT and navigate to the "Connectors" page.
2) At the bottom of the page, place a checkmark in "Advanced mode"
3) Under DMR/Homebrew/MMDVM, select "MMDVM"
4) Specify the IP address and port to connect to DMRcmd in the "Server address" and "Port (UDP)" fields respectively.
5) Click "Add Server"
6) Click "Save"
7) Update the DMRcmd config file to use the specified port.
8) Add the DMR ID (usually 9 digits) and server password into the CMDcmd config file clients section.
9) Start or re-start DMRcmd to read the new config file (by default dmrcmd.json)

Note that the "Homebrew" protocol is not supported. You must use "MMDVM" mode.

If DMRcmd is configured in proxy mode, it will proxy the server authentication as well. As a result, you must
configure OpenSPOT with the DMR ID and password that the DMR network (such as BrandMeister) expects. In proxy mode,
DMRcmd will *not* authenticate the hotspot or verify the credentials.

**In proxy mode, all DMR traffic, including to destinations intended to trigger events, will be forwarded to the 
destination DMR server unless the destination DMR ID is listed in the "drop" section of the hotspot configuration.
DMR IDs listed in the "drop" section will be dropped by the proxy in both directions, source and destination, group and private calls.**

### Home Assistant Integration

Actions can trigger Home Assistant (HA) scenes and scripts. The HA section in the configuration file must be enabled
and you must update the URL and TOKEN for your HA server.

While triggering scenes does work, scripts seem to be a more reliable integration point, especially if multiple
actions or delays are desired.

Note that by default HA uses HTTP for web console and API access. HTTP is not secure and the HA API should
therefore not be exposed to the Internet or any other untrusted network.

### License, Additional Terms, and Disclaimers

This software is licensed under GPL v3 and is intended for amateur radio and educational use only.
Please see the included LICENSE file.

If GPL presents any issues or concerns, please contact the author to discuss alternative arrangements.

By using this software you assume all risks whatsoever. This software is provided "as is"
without any warranty of any kind, either expressed or implied, including, but not limited to,
the implied warranties of merchantability and fitness for a particular purpose.

It is your responsibility to ensure that this program it not used in any way that places persons or
property at any risk.

If you do not agree to the licence and these terms, you are prohibited from using this software
or any part thereof.
