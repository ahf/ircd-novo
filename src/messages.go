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
    "container/list"
    "fmt"
    "strings"
)

func nextToken(s *string) string {
    if s == nil {
        return ""
    }

    i := strings.Index(*s, " ")

    if i == -1 {
        tmp := *s
        *s = ""
        return tmp
    }

    token := (*s)[:i]
    rest := (*s)[i + 1:]
    *s = rest

    return token
}

func ParseMessage(line string) *RawMessage {
    var prefix string

    if strings.HasPrefix(line, ":") {
        line = line[1:]
        prefix = nextToken(&line)
    }

    command := nextToken(&line)
    arguments := list.New()

    for len(line) != 0 {
        if strings.HasPrefix(line, ":") {
            arguments.PushBack(line[1:])
            break
        }

        token := nextToken(&line)
        arguments.PushBack(token)
    }

    return &RawMessage{prefix, command, arguments}
}

type RawMessage struct {
    prefix string
    command string
    arguments *list.List
}

func (r *RawMessage) String() string {
    a := "["

    for e := r.arguments.Front(); e != nil; e = e.Next() {
        value := e.Value.(string)

        if e == r.arguments.Back() {
            a = a + SingleQuote(value)
        } else {
            a = a + SingleQuote(value) + ", "
        }
    }

    a = a + "]"

    return fmt.Sprintf("Command = '%s', Parameters (Count: %d) = %s and Prefix = '%s'", r.command, r.arguments.Len(), a, r.prefix)
}

func (r *RawMessage) Command() string {
    return r.command
}

func (r *RawMessage) Prefix() string {
    return r.prefix
}

func (r *RawMessage) Arguments() *list.List {
    return r.arguments
}
