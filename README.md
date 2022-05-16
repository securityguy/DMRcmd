# DMRcmd

This program acts as an amateur radio DMR network server and executes commands based on the traffic it receives.
It was inspired by PiStar-Remote by Andy Taylor (MW0MWZ).

### Cautions

The DMR protocol is not secure and this application allows the execution of
the system commands specified in the configuration file. Anyone with a suitable receiver can
monitor DMR traffic, view the source and destination IDs, and program a radio to duplicate the
traffic. 

DMRcmd is intended for use on a private network. While it ignores traffic from unauthenticated sources,
the DMR data packets are sent in unauthenticated UDP datagrams, making it easy to forge datagrams and trigger events.

Users **must** ensure that this program it **not** used in a way that places persons or 
property at risk. It is not suitable for any application even remotely associated with life safety.
It should not be run on a computer containing sensitive or valuable information.

By default, this program does not execute any events. The user is solely responsible for what they
configure it to do.

It is up to each user to determine the suitability of this software for their specific application.

This software is intended for amateur radio and educational purposes only. 

### Compiling

This program is written in Go (aka golang). To compile:

1) If you have not already done so, install Go from https://golang.org/ or using your package manager
2) Download or clone the repo from https://github.com/securityguy/DMRcmd
3) Change to the DMRCmd directory and type `go build`

This software should compile and run on various flavours of Linux, Windows, and macOS.
If you encounter any errors please create an issue in GitHub.

Notes:
 - Once compiled, the binary and configuration file can simply be copied to another computer.
 - Go is only required to compile the software. It is not required to run the compiled program.
 - Windows users should consider installing Git from https://git-scm.com/
 
### Configuration

By default, configuration information is read from dmrcmd.conf in the working directory. 
The full path to the configuration file can optionally be passed as the first (and only) 
command line argument.

Additional configuration documentation will be added at a later date. In the interim, 
please refer to dmrcmd.conf for an example. 

For security reasons, this program neither searches the execution path nor uses a shell to 
execute commands. The full path to the program or script to be executed must be provided
in the configuration file.

For similar reasons, command line arguments to be passed to the program or script must be specified as
a list. Note that information from the DMR transmission (source ID, destination ID, client (repeater) ID, 
and IP address) can be substituted for command line arguments using $src, $dst, $ip, and $client respectively.

Within each event definition, all specified conditions must be met for the event to execute.

By design, a single DMR transmission may result in multiple events executing.

In order to avoid accidentally triggering events, a minimum number of DMR messages can be required.
The default threshold of 18 requires the transmission to last for approximately one second 
before triggering an event. Testing revealed that even a quick press and release of the PTT
sends upwards of 16 data packets, so lowering the value below 18 will likely result
in the event triggering if the PTT is bumped. If this is not a concern, the "minimum" value in the
configuration file can be changed to 1.

### Pi-Star Users

Pi-Star contains a "DRM Gateway" that is capable of routing DMR traffic to multiple servers. The most
straightforward approach to using this software is to send a range of DMR IDs to
DMRcmd, while allowing the remainder of your DMR traffic to flow to
Brandmeister and/or DMR+ as usual.

For example, in Configuration -> Expert -> Full Edit -> DMR GW:

    [DMR Network 2]
     Enabled=1
     Address=<IP where DMRcmd can be reached>
     Port=55555
     PCRewrite=2,8999900,2,8999900,100
     TGRewrite=2,8999900,2,8999900,100
     Password=PASSWORD
     Debug=0
     Id=<DMR ID>
     Name=DMRcmd

This example will send any private or group calls in the 8999900 - 8999999 range to DMRcmd.
Using a Private Call is recommended. DMRCmd ignores group calls unless "talkgroup" is
set to "true" in the event configuration, in which case both private and group calls will be
considered for the event.

The DMR ID and password in the Pi-Star DMR Gateway configuration must match a configuration in DMRcmd
or authentication will fail and no messages will be sent.

### Licence and Disclaimers

Copyright (c) 2020-2022 by Eric Jacksch VE3XEJ.

This software is licensed under GPL v3 and is intended for amateur radio and educational use only.
Please see the included LICENSE file.

By using this software you assume all risks whatsoever. This software is provided "as is"
without any warranty of any kind, either expressed or implied, including, but not limited to,
the implied warranties of merchantability and fitness for a particular purpose.
 
If you do not agree to the licence and these terms, you are prohibited from using this software
or any part thereof.
