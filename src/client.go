/* vim: set sw=4 sts=4 et foldmethod=syntax : */

/*
 * Copyright (c) 2011 Alexander Færøy <ahf@0x90.dk>
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * * Redistributions of source code must retain the above copyright notice, this
 *   list of conditions and the following disclaimer.
 *
 * * Redistributions in binary form must reproduce the above copyright notice,
 *   this list of conditions and the following disclaimer in the documentation
 *   and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package main

import (
    "bufio"
    "fmt"
    "net"
    "strconv"
    "strings"
    "time"
)

type Client struct {
    input chan *Message // Parsed IRC messages.
    output chan string // Output.

    output_quit chan bool // Leave the Output Handler.
    protocol_handler_quit chan bool // Leave the (Un)registered protocol handler.

    nickname string // The Nickname of the client.
    username string // The Username of the client.
    realname string // The Realname/Gecos of the client.
    hostname string // The hostname of the client (after we have done a reverse DNS lookup).

    address string // The IP address of the client.
    port int // The port.

    connection net.Conn // The connection.
    readwriter *bufio.ReadWriter // The buffered reader/writer.
    listener Listener // The listener.

    ircd *Ircd // Pointer to our IRCd.

    channels *ChannelSet // Channels.

    timestamp int64 // Creation time in seconds since UNIX epoch.

    // Housekeeping (FIXME: Needs refactoring).
    isRegistered bool
}

func NewClient(listener Listener, ircd *Ircd, connection net.Conn, remoteAddr string) {
    client := new(Client)

    // The IRCd.
    client.ircd = ircd

    // Timestamp.
    client.timestamp = time.Seconds()

    // IRC Channels.
    client.channels = NewChannelSet()

    // Channels.
    client.input = make(chan *Message)
    client.output = make(chan string)

    // Quit Channels.  These are being used to signal to each of our client
    // processes that it's time to exit.
    client.output_quit = make(chan bool)
    client.protocol_handler_quit = make(chan bool)

    // RemoteAddr is on the form "IP:Port". Split it.
    address, portStr, error := net.SplitHostPort(remoteAddr)

    if error != nil {
        panic("Bug: Unable to split host/port.")
    }

    // Convert it to an integer.
    port, error := strconv.Atoi(portStr)

    if error != nil {
        panic("Bug: Unable to convert port to an int.")
    }

    // Set address and port.
    client.address = address
    client.port = port

    // Start background job for looking up the hostname.
    hostnameCh := LookupAddress(address)

    // Set initial values.
    client.nickname = "*"
    client.username = "*"
    client.realname = "*"

    // Use the IP for now.
    client.hostname = address

    // Set connection, listener and create readwriter.
    client.connection = connection
    client.listener = listener
    client.readwriter = bufio.NewReadWriter(bufio.NewReader(connection), bufio.NewWriter(connection))

    // Used for synchronisation once we are done looking up rDNS and do identd
    // check.
    sync := make(chan bool)

    // Socket Output Handler.
    go client.OutputHandler()

    // Socket Input Handler.
    go client.InputHandler()

    // Unregistered Protocol Handler.  The Unregistered Protocol Handler only
    // accepts either the NICK or USER command. Once the client has been
    // registered, the UnregisteredProtocolHandler() will die and the
    // RegisteredProtocolHandler() will take over.
    go client.UnregisteredProtocolHandler(sync)

    // Start sending IRC messages to the client.
    client.WriteStringF("NOTICE AUTH :*** Processing connection to %s ...", ircd.Me())
    client.WriteString("NOTICE AUTH :*** Looking up your hostname ...")

    // Get result from rDNS lookup.
    host := <-hostnameCh
    client.WriteStringF("NOTICE AUTH :*** Found your hostname (%s) ...", host)
    client.hostname = host

    // Let the UnregisteredProtocolHandler() know that it's okay to exit and
    // let the RegisteredProtocolHandler() take over now.
    sync<-true
}

func (this *Client) InputHandler() {
    this.Printf("Entering Socket Input Handler")
    defer this.Printf("Leaving Socket Input Handler")

    rw := this.readwriter

    for {
        // Read until a newline.
        s, error := rw.ReadString('\n')

        if error != nil {
            // Read Error (Usually happens when the connection is closed by the
            // client).
            // FIXME: Should probably check for non-EOF errors.
            this.Printf("Read Error: %s", error)

            // Start a chain of Quit-messages to the various client processes
            // to ensure that they are all shutdown gracefully. The goal of
            // this is to avoid any dangling processes.
            this.output_quit<-true

            // Exit.
            return
        }

        // Trim "\r\n".
        line := strings.TrimRight(s, "\r\n")

        // Skip "empty" lines.
        if len(Strip(line)) == 0 {
            continue
        }

        // Parse.
        message := Parse(line)

        // Most ASCII strings are valid IRC commands, but in the
        // unlikely case that Parse() returns nil, we'll simply skip
        // it.
        if message != nil {
            this.Printf("-> '%s'", message)

            // Pass on to the ProtocolHandler (Either Registered or Unregistered).
            this.input<-message
        }
    }
}

func (this *Client) OutputHandler() {
    this.Printf("Entering Socket Output Handler")
    defer this.Printf("Leaving Socket Output Handler")

    defer func () {
        this.Printf("Sending Quit-message to the Protocol Handler")
        this.protocol_handler_quit<-true
    }()

    rw := this.readwriter

    for {
        select {
            case <-this.output_quit:
                // We were told to exit. Notice the above defer statements
                // where we chain the quit message on to the message handler
                // process.
                return

            case s, ok := <-this.output:
                if ! ok {
                    // The client is dead.
                    // FIXME: Shouldn't happen.
                    return
                }

                // Ignore the byte-count for now.  Go sets error to non-nil if
                // bytes != len(input).
                _, error := rw.WriteString(s + "\r\n")

                if error != nil {
                    this.Printf("Write Error: %s", error)

                    // Unlikely, but we'll quit.
                    return
                }

                this.Printf("<- '%s'", s)

                // Flush the buffer.
                rw.Flush()
        }
    }
}

func (this *Client) UnregisteredProtocolHandler(sync chan bool) {
    this.Printf("Entering Unregistered Client Protocol Handler")
    defer this.Printf("Leaving Unregistered Client Protocol Handler")

    // The IRCd.
    ircd := this.ircd

    // Housekeeping.
    got_nickname := false
    got_username := false

    // Temporaries.
    var nickname string
    var username string
    var realname string

    for {
        select {
            case message, ok := <-this.input:
                if ! ok {
                    // This client is dead.
                    // FIXME: Shouldn't happen.
                    return
                }

                cmd := message.Command()
                args := message.Arguments()

                switch cmd {
                    case "NICK":
                        // Missing nickname.
                        if len(args) < 1 {
                            this.SendNumeric(ERR_NONICKNAMEGIVEN, ircd.Me(), "*")
                            continue
                        }

                        n := args[0]

                        // Check if the nickname is valid.
                        if ! IsValidNickname(n) {
                            this.SendNumeric(ERR_ERRONEUSNICKNAME, ircd.Me(), "*", n)
                            continue
                        }

                        // Check if the nickname is already registered.
                        if this.ircd.FindClient(n) != nil {
                            this.SendNumeric(ERR_NICKNAMEINUSE, ircd.Me(), "*", n)
                            continue
                        }

                        got_nickname = true
                        nickname = n

                    case "USER":
                        // Missing parameters.
                        if len(args) < 4 {
                            this.SendNumeric(ERR_NEEDMOREPARAMS, ircd.Me(), "*", cmd)
                            continue
                        }

                        u := args[0]
                        r := args[3]

                        // Check if the username is valid.
                        if ! IsValidUsername(u) {
                            // FIXME: Handle ...
                        }

                        got_username = true

                        // FIXME: We Check Identd later.
                        username = "~" + u
                        realname = r
                }

                if got_nickname && got_username {
                    // Update the client information.  NOTE: It's important to
                    // do this *before* any attempts to register using
                    // this.Register() since otherwise it'll try to register
                    // using the "default" nickname, which is "*", and weird
                    // things may or may not happen.
                    this.nickname = nickname
                    this.username = username
                    this.realname = realname

                    // Try doing a registration.
                    // Check if the nickname has already been registered, just
                    // in case someone else connected whilst we were connecting
                    // and stole our nickname.
                    if ! this.Register() {
                        this.SendNumeric(ERR_NICKNAMEINUSE, ircd.Me(), "*", nickname)
                        continue
                    }

                    // Wait for the synchronisation signal.  We have to do rDNS
                    // lookup and identd check before we proceed to the
                    // registered protocol handler.
                    this.Printf("Waiting for synchronisation")
                    <-sync
                    this.Printf("Synchronisation done")

                    // Let the full IRC protocol handler take over now.
                    go this.RegisteredProtocolHandler()
                    return
                }

            case <-this.protocol_handler_quit:
                // We were told to quit.
                return
        }
    }
}

func (this *Client) RegisteredProtocolHandler() {
    this.Printf("Entering Registered Client Protocol Handler!")
    defer this.Printf("Leaving Registered Client Protocol Handler!")

    // The IRCd.
    ircd := this.ircd

    // The server information string.
    server := fmt.Sprintf("%s[%s]", ircd.Me(), this.listener.Address())

    // Welcome the client.
    this.SendNumeric(RPL_WELCOME, ircd.Me(), this.Nickname(), ircd.Description(), this.Nickname())
    this.SendNumeric(RPL_YOURHOST, ircd.Me(), this.Nickname(), server, VersionFull)

    // Send MOTD, if any.
    this.SendMotd()

    for {
        select {
            case message, ok := <-this.input:
                if ! ok {
                    // The client is dead.
                    // FIXME: Shouldn't happen.
                    return
                }

                // Handle the message.
                HandleMessage(this, message)

            case <-this.protocol_handler_quit:
                // We were told to quit.
                return
        }
    }
}

func (this *Client) WriteString(s string) {
    this.output<-s
}

func (this *Client) WriteStringF(format string, a...interface{}) {
    this.WriteString(fmt.Sprintf(format, a...))
}

func (this *Client) SendNumeric(n Numeric, a...interface{}) {
    this.WriteString(Format(n, a...))
}

func (this *Client) Unregister() {
    this.Printf("Unregistering client")

    if this.isRegistered {
        this.ircd.UnregisterClient(this)
    }
}

func (this *Client) Register() bool {
    this.Printf("Trying to register client")

    if this.isRegistered {
        panic("Bug: Trying to register an already registered client.")
    }

    b := this.ircd.RegisterClient(this)

    this.isRegistered = b
    return b
}

func (this *Client) SendMotd() {
    ircd := this.ircd
    motd := ircd.Motd()

    if motd == nil {
        this.SendNumeric(ERR_NOMOTD, ircd.Me(), this.Nickname())
        return
    }

    lines := *motd

    this.SendNumeric(RPL_MOTDSTART, ircd.Me(), this.Nickname(), ircd.Me())

    for i := range lines {
        this.SendNumeric(RPL_MOTD, ircd.Me(), this.Nickname(), lines[i])
    }

    this.SendNumeric(RPL_ENDOFMOTD, ircd.Me(), this.Nickname())
}

func (this *Client) Quit(message string) {
    // FIXME: Hack: Should be fixed like the description in bug #13.
    // Send a PART for now.
    this.channels.ForEach(func (channel *Channel) {
        channel.Part(this)
    })

    this.WriteStringF(":%s QUIT :%s", this, message)

    // The Socket Input Handler will start the chain that kills the other
    // processes.
    this.Close()
}

func (this *Client) Pong(message string) {
    this.WriteStringF("PONG :%s", message)
}

func (this *Client) Close() {
    this.Printf("Closing")
    defer this.Printf("Finished Closing")

    // Unregister the client.
    this.Unregister()

    // Cleanup.
    this.connection.Close()
}

func (this *Client) PrivateMessage(message *PrivateMessage) {
    this.WriteString(message.String())
}

func (this *Client) Printf(format string, a...interface{}) {
    this.ircd.Printf(this.String() + ": " + format, a...)
}

func (this *Client) RemotePort() int {
    return this.port
}

func (this *Client) RemoteAddress() string {
    return this.address
}

func (this *Client) Nickname() string {
    return this.nickname
}

func (this *Client) Username() string {
    return this.username
}

func (this *Client) Hostname() string {
    return this.hostname
}

func (this *Client) Realname() string {
    return this.realname
}

func (this *Client) Join(channel *Channel) {
    channel.Join(this)
    this.channels.Insert(channel)
}

func (this *Client) Part(channel *Channel) {
    channel.Part(this)
    this.channels.Delete(channel)
}

func (this *Client) ChannelJoin(source *Client, channel *Channel) {
    this.WriteStringF(":%s JOIN :%s", source, channel)
}

func (this *Client) ChannelPart(source *Client, channel *Channel) {
    this.WriteStringF(":%s PART :%s", source, channel)
}

// FIXME: Should die in a refactoring session. Only needed by the
// CommandRegistry.
func (this *Client) Ircd() *Ircd {
    return this.ircd
}

func (this *Client) String() string {
    return fmt.Sprintf("%s!%s@%s", this.Nickname(), this.Username(), this.Hostname())
}
