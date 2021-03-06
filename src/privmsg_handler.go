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

func init() {
    RegisterMessageHandler("PRIVMSG", PrivmsgHandler)
}

func PrivmsgHandler(source *Client, message *Message) {
    cmd := message.Command()
    args := message.Arguments()
    count := len(args)
    ircd := source.Ircd()

    if count < 1 {
        source.SendNumeric(ERR_NORECIPIENT, ircd.Me(), source.Nick(), cmd)
        return
    }

    if count < 2 {
        source.SendNumeric(ERR_NOTEXTTOSEND, ircd.Me(), source.Nick())
        return
    }

    target := args[0]
    text := args[1]

    // Create our message.
    msg := NewPrivateMessage(source, target, text)

    // Try to find a client.
    client := ircd.FindClient(target)

    if client != nil {
        // client.WriteStringF(":%s PRIVMSG %s :%s", source, n.Nick(), text)
        client.PrivateMessage(msg)
        return
    }

    // No client found? Look for a channel.
    channel := ircd.FindChannel(target)

    if channel != nil {
        channel.PrivateMessage(msg)
        return
    }

    // Error out:
    // FIXME: No Such Nick/Chan...
}
