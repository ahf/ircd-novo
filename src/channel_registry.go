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

type ChannelRegistry struct {
    ircd *Ircd // Pointer to the IRCd instance.

    // Map channel names to Channels
    channels map[string] *Channel
}

func NewChannelRegistry(ircd *Ircd) *ChannelRegistry {
    return &ChannelRegistry{ircd, make(map[string] *Channel)}
}

func (this *ChannelRegistry) Find(name string) *Channel {
    if channel, ok := this.channels[ToLower(name)]; ok {
        return channel
    }

    return nil
}

func (this *ChannelRegistry) FindOrCreate(name string) *Channel {
    lowered := ToLower(name)
    channel := this.Find(lowered)

    if channel == nil {
        // Don't use ToLower() here.
        channel = NewChannel(this.ircd, name)

        // But use it here.
        this.channels[lowered] = channel
    }

    return channel
}

func (this *ChannelRegistry) Unregister(channel *Channel) {
    this.Printf("Unregistering Channel", channel)
    this.channels[ToLower(channel.Name())] = nil, false
}

func (this *ChannelRegistry) Printf(format string, a...interface{}) {
    this.ircd.Printf(format, a...)
}
