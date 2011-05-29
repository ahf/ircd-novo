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

type ClientSet struct {
    // FIXME: Due to Go's lack of build-in type-safe set's, we are currently
    // going to (ab)use a map of *Client's to bool's to represent a set. The
    // boolean value is simply ignored for all operations.
    clients map[*Client] bool
}

func NewClientSet() *ClientSet {
    return &ClientSet{make(map[*Client] bool)}
}

func (this *ClientSet) Insert(client *Client) {
    this.clients[client] = true
}

func (this *ClientSet) Contains(client *Client) bool {
    _, exists := this.clients[client]
    return exists
}

func (this *ClientSet) Delete(client *Client) {
    this.clients[client] = false, false
}

func (this *ClientSet) ForEach(f func (*Client)) {
    for client, _ := range this.clients {
        f(client)
    }
}

func (this *ClientSet) Len() int {
    return len(this.clients)
}

func (this *ClientSet) Names() []string {
    r := make([]string, len(this.clients))
    i := 0

    for client := range this.clients {
        r[i] = client.Nick()
        i++
    }

    return r
}
