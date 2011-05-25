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
    "strings"
)

func init() {
    RegisterMessageHandler("WHOIS", WhoisHandler)
}

func WhoisHandler(client *Client, message *Message) {
    args := message.Arguments()
    ircd := client.Ircd()

    if len(args) < 1 {
        client.SendNumeric(ERR_NEEDMOREPARAMS, client.Ircd().Me(), client.Nickname(), message.Command())
        return
    }

    target := ircd.FindClient(args[0])

    if target == nil {
        // FIXME: No such nick/chan
        return
    }

    // FIXME: Handle local whois too.

    // Send info.
    client.SendNumeric(RPL_WHOISUSER, ircd.Me(), client.Nickname(), target.Nickname(), target.Username(), target.Hostname(), target.Realname())
    client.SendNumeric(RPL_WHOISCHANNELS, ircd.Me(), client.Nickname(), target.Nickname(), strings.Join(target.ChannelNames(), " "))
    client.SendNumeric(RPL_WHOISSERVER, ircd.Me(), client.Nickname(), target.Nickname(), ircd.Me(), ircd.Description())

    // if client is an operator {
    //     client.SendNumeric(RPL_WHOISOPERATOR)
    // }

    if client.Secure() {
        client.SendNumeric(RPL_WHOISSECURE, ircd.Me(), client.Nickname(), target.Nickname())
    }

    if client.WebSocket() {
        client.SendNumeric(RPL_WHOISWEBSOCKET, ircd.Me(), client.Nickname(), target.Nickname())
    }

    // FIXME: Idle Times.
    // client.SendNumeric(RPL_WHOISIDLE)

    client.SendNumeric(RPL_ENDOFWHOIS, ircd.Me(), client.Nickname(), target.Nickname())
}
