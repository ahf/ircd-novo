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

type ClientRegistry struct {
    ircd *Ircd

    // Maps nicknames to a Client.
    clients map[string] *Client
}

func NewClientRegistry(ircd *Ircd) *ClientRegistry {
    return &ClientRegistry{ircd, make(map[string] *Client)}
}

func (this *ClientRegistry) Find(nick string) *Client {
    if client, ok := this.clients[ToLower(nick)]; ok {
        return client
    }

    return nil
}

func (this *ClientRegistry) Register(client *Client) bool {
    this.Printf("Registering %s", client.Nickname())

    n := ToLower(client.Nickname())

    if _, exists := this.clients[n]; exists {
        return false
    }

    this.clients[n] = client
    return true
}

func (this *ClientRegistry) Unregister(client *Client) {
    this.Printf("Unregistering %s", client)
    this.clients[ToLower(client.Nickname())] = nil, false
}

func (this *ClientRegistry) Printf(format string, a...interface{}) {
    this.ircd.Printf(format, a...)
}
