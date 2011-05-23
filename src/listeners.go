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
    "crypto/tls"
    "fmt"
    "http"
    "log"
    "net"
    "os"
    "strings"
    "websocket"
)

type Protocol int

const (
    TCP Protocol = iota
    WebSocket
)

func ProtocolFromString(p string) *Protocol {
    switch strings.ToLower(p) {
        case "websocket":
            w := WebSocket
            return &w
        case "tcp":
            t := TCP
            return &t
    }

    return nil
}

func (this Protocol) String() string {
    switch this {
        case TCP: return "TCP"
        case WebSocket: return "WebSocket"
    }

    panic("Unhandled Protocol.")
}

type SecurityLevel int

const (
    Insecure SecurityLevel = iota
    Secure
)

func (this SecurityLevel) String() string {
    switch this {
        case Insecure: return "Insecure"
        case Secure: return "Secure"
    }

    panic("Unhandled SecurityLevel.")
}

type Listener interface {
    Protocol() Protocol
    Address() net.Addr
    Secure() SecurityLevel
    Listen()
}

type BasicListener struct {
    *log.Logger
    secure SecurityLevel
    protocol Protocol
    listener net.Listener
    ircd *Ircd
}

func NewBasicListener(ircd *Ircd) *BasicListener {
    logger := log.New(os.Stderr, "", log.Ldate | log.Ltime)

    return &BasicListener{logger, Insecure, TCP, nil, ircd}
}

func (this *BasicListener) Protocol() Protocol {
    return this.protocol
}

func (this *BasicListener) Address() net.Addr {
    return this.listener.Addr()
}

func (this *BasicListener) Secure() SecurityLevel {
    return this.secure
}

func (this *BasicListener) Listen() {
    panic("Listeners must implement Listen() themselves.")
}

func (this *BasicListener) ConnectionHandler(ircd *Ircd, connection net.Conn, remoteAddr string) {
    this.Printf("Incoming Connection from %s to %s (%s %s)", remoteAddr, this.Address(), this.Secure(), this.Protocol())
    NewClient(ircd, connection, remoteAddr)
}

type TCPListener struct {
    *BasicListener
}

func NewTCPListener(ircd *Ircd, address string, port int, config *tls.Config) (*TCPListener, os.Error) {
    var listener net.Listener
    secure := Insecure

    // Go is slightly silly here. The IPv6 form, in Go, is "[::1]:6667" whereas
    // the IPv4 form is "127.0.0.1:6667". Let's try to satisfy both parties by
    // firstly trying with 'a:6667' and if that fails retry with '[a]:6667'.
    finalAddress, error := net.ResolveTCPAddr("TCP", address)

    if error != nil {
        finalAddress, error = net.ResolveTCPAddr("TCP", fmt.Sprintf("[%s]:%d", address, port))

        if error != nil {
            return nil, error
        }
    }

    listener, error = net.ListenTCP(finalAddress.Network(), finalAddress)

    if error != nil {
        return nil, error
    }

    if config != nil {
        listener = tls.NewListener(listener, config)
        secure = Secure
    }

    bl := NewBasicListener(ircd)
    bl.listener = listener
    bl.secure = secure

    return &TCPListener{bl}, nil
}

func (this *TCPListener) Listen() {
    for {
        connection, error := this.listener.Accept()

        if error != nil {
            this.Printf("Error: %s", error)
            continue
        }

        // Set the RemoteAddr here because of Go Bug #1636
        go this.ConnectionHandler(this.ircd, connection, connection.RemoteAddr().String())
    }
}

type WebSocketListener struct {
    *TCPListener
}

func NewWebSocketListener(ircd *Ircd, address string, port int, config *tls.Config) (*WebSocketListener, os.Error) {
    t, error := NewTCPListener(ircd, address, port, config)

    if error != nil {
        return nil, error
    }

    wsl := &WebSocketListener{t}
    wsl.protocol = WebSocket
    return wsl, nil
}

func (this *WebSocketListener) Listen() {
    http.Serve(this.listener, websocket.Handler(func (connection *websocket.Conn) {
        // The HTTP package creates the goroutine itself. No need for us to do it.
        // Set the RemoteAddr here because of Go Bug #1636
        this.ConnectionHandler(this.ircd, connection, connection.Request.RemoteAddr)
    }))
}
