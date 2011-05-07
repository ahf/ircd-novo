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
    "log"
    "os"
    "time"
)
type Ircd struct {
    *log.Logger
    listeners []Listener

    certificateFile *string
    keyFile *string
}

func NewIrcd() *Ircd {
    return &Ircd{log.New(os.Stderr, "", log.Ldate | log.Ltime), make([]Listener, 0), nil, nil}
}

func (this *Ircd) SetCertificate(certificate, key string) {
    this.certificateFile = &certificate
    this.keyFile = &key
}

func (this *Ircd) addListener(p Protocol, address string, config *tls.Config) {
    var listener Listener

    switch p {
        case TCP: listener = NewTCPListener(address, config)
        case WebSocket: listener = NewWebSocketListener(address, config)
        default: panic("Unhandled Protocol.")
    }

    if listener != nil {
        this.listeners = append(this.listeners, listener)
    }
}

func (this *Ircd) AddListener(protocol Protocol, address string) {
    this.addListener(protocol, address, nil)
}

func (this *Ircd) AddSecureListener(protocol Protocol, address string) {
    if this.certificateFile == nil || this.keyFile == nil {
        this.Printf("Unable to add secure listener. Missing certificate or key file")
        return
    }

    cert, error := tls.LoadX509KeyPair(*this.certificateFile, *this.keyFile)

    if error != nil {
        this.Printf("Error Loading Certificate: %s", error)
        return
    }

    config := &tls.Config{
        Rand: rand.Reader,
        Time: time.Seconds,
    }

    config.Certificates = make([]tls.Certificate, 1)
    config.Certificates[0] = cert

    if protocol == WebSocket {
        config.NextProtos = []string{"http/1.1"}
    }

    this.addListener(protocol, address, config)
}

func (this *Ircd) Listen() {
    for i := range this.listeners {
        listener := this.listeners[i]

        this.Printf("Listening on %s (%s %s)", listener.Address(), listener.Secure(), listener.Protocol())
        go listener.Listen()
    }
}
