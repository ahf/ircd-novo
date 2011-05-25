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
    "crypto/rand"
    "crypto/tls"
    "fmt"
    "io/ioutil"
    "log"
    "os"
    "strings"
    "time"
)
type Ircd struct {
    *log.Logger
    listeners []Listener
    config *ConfigurationFile

    motdFile string
    motdContent []string

    clientRegistry *ClientRegistry
    channelRegistry *ChannelRegistry
}

func NewIrcd() *Ircd {
    ircd := new(Ircd)
    ircd.Logger = log.New(os.Stderr, "", log.Ldate | log.Ltime)
    ircd.listeners = make([]Listener, 0)
    ircd.clientRegistry = NewClientRegistry(ircd)
    ircd.channelRegistry = NewChannelRegistry(ircd)

    return ircd
}

func (this *Ircd) SetConfigurationFile(config *ConfigurationFile) {
    this.config = config

    for i := range this.config.Ircd.Listeners {
        listener := this.config.Ircd.Listeners[i]
        protocol := ProtocolFromString(listener.Type)
        host := listener.Host
        port := listener.Port

        if protocol == nil {
            this.Printf("Unknown protocol type: %s\n", listener.Type)
            continue
        }

        if listener.Tls {
            this.addSecureListener(*protocol, host, port)
        } else {
            this.addListener(*protocol, host, port)
        }
    }
}

func (this *Ircd) addCommonListener(p Protocol, address string, port int, config *tls.Config) {
    var listener Listener
    var error os.Error

    switch p {
        case TCP: listener, error = NewTCPListener(this, address, port, config)
        case WebSocket: listener, error = NewWebSocketListener(this, address, port, config)
        default: panic("Unhandled Protocol.")
    }

    if error != nil {
        this.Printf("Error: %s", error)
        return
    }

    this.listeners = append(this.listeners, listener)
}

func (this *Ircd) addListener(protocol Protocol, address string, port int) {
    this.addCommonListener(protocol, address, port, nil)
}

func (this *Ircd) addSecureListener(protocol Protocol, address string, port int) {
    cert := this.config.Ircd.ServerInfo.Tls.Certificate
    key := this.config.Ircd.ServerInfo.Tls.Key
    errorMessage := fmt.Sprintf("Unable to add secure listener for %s:%d", address, port)

    if cert == "" {
        this.Printf("%s: %s", errorMessage, "Empty TLS certificate in configuration file.")
        return
    }

    if key == "" {
        this.Printf("%s: %s", errorMessage, "Empty TLS key in configuration file.")
        return
    }

    certificate, error := tls.LoadX509KeyPair(cert, key)

    if error != nil {
        this.Printf("Error Loading Certificate: %s", error)
        return
    }

    config := &tls.Config{
        Rand: rand.Reader,
        Time: time.Seconds,
    }

    config.Certificates = make([]tls.Certificate, 1)
    config.Certificates[0] = certificate

    if protocol == WebSocket {
        config.NextProtos = []string{"http/1.1"}
    }

    this.addCommonListener(protocol, address, port, config)
}

func (this *Ircd) Run() {
    if len(this.listeners) == 0 {
        fmt.Printf("Error: No Listeners Defined...\n")
        os.Exit(1)
    }

    this.Printf("Opening up for incoming connections")

    for i := range this.listeners {
        listener := this.listeners[i]

        this.Printf("Listening on %s (%s %s)", listener.Address(), listener.Secure(), listener.Protocol())
        go listener.Listen()
    }
}

func (this *Ircd) Me() string {
    return this.config.Ircd.ServerInfo.Name
}

func (this *Ircd) Description() string {
    return this.config.Ircd.ServerInfo.Description
}

func (this *Ircd) SetMotdFile(path string) {
    this.motdFile = path

    this.LoadMotd()
}

func (this *Ircd) LoadMotd() {
    content, error := ioutil.ReadFile(this.motdFile)

    if error != nil {
        this.Printf("Unable to load MOTD file: %s", error)
        return
    }

    this.motdContent = strings.Split(string(content), "\n", -1)
}

func (this *Ircd) Motd() *[]string {
    if len(this.motdContent) == 0 {
        return nil
    }

    return &this.motdContent
}

func (this *Ircd) FindClient(nick string) *Client {
    return this.clientRegistry.Find(nick)
}

func (this *Ircd) RegisterClient(client *Client) bool {
    return this.clientRegistry.Register(client)
}

func (this *Ircd) UnregisterClient(client *Client) {
    this.clientRegistry.Unregister(client)
}

func (this *Ircd) FindChannel(name string) *Channel {
    return this.channelRegistry.Find(name)
}

func (this *Ircd) FindOrCreateChannel(name string) *Channel {
    return this.channelRegistry.FindOrCreate(name)
}

func (this *Ircd) UnregisterChannel(channel *Channel) {
    this.channelRegistry.Unregister(channel)
}

func (this *Ircd) ForEachChannel(f func (*Channel)) {
    this.channelRegistry.ForEach(f)
}
