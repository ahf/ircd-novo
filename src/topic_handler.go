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
    RegisterMessageHandler("TOPIC", TopicHandler)
}

func TopicHandler(client *Client, message *Message) {
    args := message.Arguments()
    ircd := client.Ircd()

    // Need at least the channel parameter.
    if len(args) < 1 {
        client.SendNumeric(ERR_NEEDMOREPARAMS, ircd.Me(), client.Nickname(), message.Command())
        return
    }

    // Find Channel.
    channel := ircd.FindChannel(args[0])

    // FIXME: No such channel or nick.
    if channel == nil {
        return
    }

    // Read Channel Topic.
    if len(args) < 2 {
        c := make(chan *Topic, 1)
        channel.Topic(c)
        topic := <-c

        if topic == nil {
            client.SendNumeric(RPL_NOTOPIC, ircd.Me(), client.Nickname(), channel)
            return
        }

        client.SendNumeric(RPL_TOPIC, ircd.Me(), client.Nickname(), channel, topic)
        return
    }

    text := args[1]
    channel.SetTopic(client, text)
}
