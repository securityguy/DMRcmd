# DMRcmd

Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ.

This is early beta software. Please read this file in its entirety and use at your own risk.

This program acts as an amateur radio DMR network server using the Homebrew repeater protocol (MMDVMHost variant) and
can be configured to execute commands based on the traffic it receives from the user's hotspot.

It was inspired by PiStar-Remote by Andy Taylor (MW0MWZ) and DMRGateway by Jonathan Naylor (G4KLX).

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
3) Change to the DMRCmd directory and type `go build`.

This software should compile and run on various flavours of Linux, Windows, and macOS.
If you encounter any errors please create an issue in GitHub.

Notes:
 - Once compiled, the binary and configuration file can simply be copied to another computer.
 - Go is only required to compile the software. It is not required to run the compiled program.
 - Windows users should consider installing Git from https://git-scm.com/.
 
### Configuration

By default, configuration information is read from dmrcmd.json in the working directory. 
The full path to the configuration file can optionally be passed as the first (and only) 
command line argument.

Note that most configuration objects have an "enabled" option. Unless it is set to true, the
object will be ignored.

Additional configuration documentation will be added at a later date. In the interim, 
please refer to dmrcmd.example.json.

For security reasons, this program neither searches the execution path nor uses a shell to 
execute commands. The full path to a program or script to be executed must be provided
in the configuration file.

For similar reasons, command line arguments to be passed to the program or script must be specified as
a list. Note that information from the DMR transmission (source ID, destination ID, client (repeater) ID, 
and IP address) can be substituted for command line arguments using $src, $dst, $ip, and $client respectively.

Within each event definition, **all specified conditions** must be met for the event to execute.

By design, a single DMR transmission may result in multiple events executing.

In order to avoid accidentally triggering events, a minimum number of DMR messages can be required.
The default threshold of 18 requires the transmission to last for approximately one second 
before triggering an event. Testing revealed that even a quick press and release of the PTT
sends upwards of 16 data packets, so lowering the value below 18 will likely result
in the event triggering if the PTT is bumped. If this is not a concern, the "minimum" value in the
configuration file can be changed to 1.

### Pi-Star and DRMGateway Users

Pi-Star contains a "DRM Gateway" that is capable of routing DMR traffic to multiple servers. 
DMRGateway by Jonathan Naylor (G4KLK) is also available as a stand-alone application from 
https://github.com/g4klx/DMRGateway.

The most straightforward approach to using this software is to send a range of DMR IDs to
DMRcmd, while allowing the remainder of your DMR traffic to flow to Brandmeister, DMR+, or other service as usual.
In that case, do not configure any server objects in the DMRcmd configuration file.

For example, in the DMRGateway configuration file (Pi-Star users can access it at Configuration -> Expert -> Full Edit: DMR GW):

    [DMR Network 2]
     Enabled=1
     Address=<IP where DMRcmd can be reached>
     Port=55555
     PCRewrite=2,8999900,2,8999900,100
     TGRewrite=2,8999900,2,8999900,100
     Password=<PASSWORD>
     Debug=0
     Id=<DMR ID>
     Name=DMRcmd

This example will send any private or group calls in the 8999900 - 8999999 range to DMRcmd.
Using a Private Call is recommended. DMRCmd ignores group calls unless "talkgroup" is
set to true in the event configuration, in which case both private and group calls will be
considered for the event.

The DMR ID and password in the DMR Gateway configuration must match a configuration in DMRcmd
or authentication will fail and data from your hotspot will be ignored.

**Please be cognizant of the fact that any given DMR ID may be assigned to an individual and use appropriate
care to ensure that you do not inadvertently route traffic intended to trigger DMRCmd events to a DMR network.**

### openSPOT Users

openSPOT does not include multi-DMR server routing capability, so at this point configuring openSPOT to connect to 
DMRcmd will not allow other traffic to reach a DMR server. I am in the process of adding passthrough capability. 
However, at this time your best bet would be to configure your openSPOT to send traffic to
DMRGateway (link and configuration above) and configure DMRGateay to send selected traffic to DMRcmd.

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

Also note that the procedure above adds the IP address and port you have specified to the "Custom servers" section
of the server list, so you are free to uncheck advanced mode.

### Home Assistant Integration

Actions can trigger Home Assistant scenes and scripts. The HA section in the configuration file must be enabled
and you must update the URL and TOKEN for your Home Assistant server. 

While trigger scenes does work, scripts seem to be a more reliable integration point, especially if multiple 
actions or delays are desired.

Note that by default Home Assistant uses HTTP for web console and API access. HTTP is not secure.

### Licence, Additional Terms, and Disclaimers

This software is licensed under GPL v3 and is intended for amateur radio and educational use only.
Please see the included LICENSE file.

If GPL presents any issues or concerns please contact the author to discuss alternative arrangements.

By using this software you assume all risks whatsoever. This software is provided "as is"
without any warranty of any kind, either expressed or implied, including, but not limited to,
the implied warranties of merchantability and fitness for a particular purpose.

It is your responsibility to ensure that this program it not used in any way that places persons or
property at any risk.
 
If you do not agree to the licence and these terms, you are prohibited from using this software
or any part thereof.
