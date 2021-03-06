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
    "time"
)

type Channel struct {
    ircd *Ircd // Pointer to the IRCd instance.

    name string // The name of the channel.
    topic *Topic // The topic of the channel.

    joining chan *Client // Channel of Joining Members.
    parting chan *Client // Channel of Members whom are leaving.
    read_topic chan chan *Topic // Channel for synchronous read of the topic.
    read_client_count chan chan int // Channel for reading the client count of the channel.
    private_messages chan *PrivateMessage // Channel of private messages.
    notice_messages chan *NoticeMessage // Channel of notice messages.
    write_topic chan *TopicMessage // Channel to set topic's.

    clients *ClientSet // Client Members.

    timestamp int64 // Creation time in seconds since UNIX epoch.
}

func NewChannel(ircd *Ircd, name string) *Channel {
    channel := &Channel {
        name: name, // Channel Name.
        timestamp: time.Seconds(), // Timestamp.
        topic: nil, // No topic.
        ircd: ircd, // The IRCd.
        clients: NewClientSet(), // Our client members.

        joining: make(chan *Client),
        parting: make(chan *Client),
        private_messages: make(chan *PrivateMessage),
        notice_messages: make(chan *NoticeMessage),

        read_topic: make(chan chan *Topic),
        read_client_count: make(chan chan int),

        write_topic: make(chan *TopicMessage),
    }

    // Message Handler.
    go channel.Handler()

    return channel
}

func (this *Channel) Handler() {
    this.Printf("Starting Channel Handler")
    defer this.Printf("Leaving Channel Handler")

    // The IRCd.
    ircd := this.ircd

    for {
        select {
            case joining_client := <-this.joining:
                this.Printf("Client '%s' joined.", joining_client)

                // This is true, if our joining client is the creator of the channel.
                //   creator := this.clients.Len() == 0

                // Insert our new client.
                this.clients.Insert(joining_client)

                // Send JOIN message to all clients.
                this.clients.ForEach(func (client *Client) {
                    client.ChannelJoin(joining_client, this)
                })

                // Topic, if any.
                if this.topic != nil {
                    joining_client.SendNumeric(RPL_TOPIC, ircd.Me(), joining_client.Nick(), this.name, this.topic)
                    joining_client.SendNumeric(RPL_TOPICWHOTIME, ircd.Me(), joining_client.Nick(), this.name, this.topic.Setter(), this.topic.Timestamp())
                }

                // Client Names.
                names := this.clients.Names()

                // NOTE: See RB codebase for information about the "=" here.
                joining_client.SendNumeric(RPL_NAMREPLY, ircd.Me(), joining_client.Nick(), "=", this.name, strings.Join(names, " "))

                // FIXME: This could become a long message for large channels.
                joining_client.SendNumeric(RPL_ENDOFNAMES, ircd.Me(), joining_client.Nick(), this.name)


            case parting_client := <-this.parting:
                this.Printf("Client '%s' left.", parting_client)

                // Send PART message to all clients, including ourself.
                this.clients.ForEach(func (client *Client) {
                    client.ChannelPart(parting_client, this)
                })

                // Remove our client.
                this.clients.Delete(parting_client)

                // Last member left?
                if this.clients.Len() == 0 {
                    // Unregister.
                    this.Unregister()

                    // Shutdown.
                    return
                }

            case topic_reader := <-this.read_topic:
                // Send topic.
                topic_reader<-this.topic

            case client_count_reader := <-this.read_client_count:
                // Send the client count.
                client_count_reader<-this.clients.Len()

            case topic_message := <-this.write_topic:
                // Text.
                text := topic_message.Text()

                // Source Client.
                source := topic_message.Source()

                // Set topic.
                this.topic = NewTopic(source, text)

                // Send update to all members.
                this.clients.ForEach(func (client *Client) {
                    // FIXME: WriteStringF ...
                    client.WriteStringF(":%s TOPIC %s :%s", source, this.name, text)
                })

            case message := <-this.private_messages:
                // Source Client.
                source := message.Source()

                // Broadcast to each client, except the source.
                this.clients.ForEach(func (client *Client) {
                    // FIXME: Should we compare source.Nick() with
                    // client.Nick() here?
                    if source == client {
                        // Don't send message to ourself.
                        return
                    }

                    // Send message.
                    client.PrivateMessage(message)
                })

            case message := <-this.notice_messages:
                // Source client.
                source := message.Source()

                // Broadcast to everyone, except the source.
                this.clients.ForEach(func (client *Client) {
                    if source == client {
                        // Source.
                        return
                    }

                    // Send message.
                    client.Notice(message)
                })
        }
    }
}

func (this *Channel) Join(client *Client) {
    this.joining<-client
}

func (this *Channel) Part(client *Client) {
    this.parting<-client
}

func (this *Channel) Topic(c chan *Topic) {
    this.read_topic<-c
}

func (this *Channel) SetTopic(client *Client, text string) {
    this.write_topic<-NewTopicMessage(client, text)
}

func (this *Channel) ClientCount(c chan int) {
    this.read_client_count<-c
}

func (this *Channel) PrivateMessage(message *PrivateMessage) {
    this.private_messages<-message
}

func (this *Channel) Notice(message *NoticeMessage) {
    this.notice_messages<-message
}

func (this *Channel) Unregister() {
    if this.clients.Len() != 0 {
        panic("Bug: Trying to unregister a non-empty channel.")
    }

    this.ircd.UnregisterChannel(this)
}

func (this *Channel) Name() string {
    // NOTE: This is safe to return without asking our channel handler process,
    // since a channel name can never be changed and is set upon construction
    // in NewChannel().
    return this.name
}

func (this *Channel) String() string {
    // NOTE: See information in Name().
    return this.name
}

func (this *Channel) Printf(format string, a...interface{}) {
    this.ircd.Printf(this.String() + ": " + format, a...)
}
