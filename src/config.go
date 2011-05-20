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
    "json"
    "os"
)

type ConfigurationFile struct {
    Ircd IrcdConfigType
}

type IrcdConfigType struct {
    ServerInfo ServerInfoConfigType
    NetworkInfo NetworkInfoConfigType
    AdminInfo AdminInfoConfigType
    Listeners []ListenerConfigType
    Operators []OperatorConfigType
}

type ServerInfoConfigType struct {
    Name string
    Description string
    Tls TlsConfigType
}

type TlsConfigType struct {
    Certificate string
    Key string
}

type NetworkInfoConfigType struct {
    Name string
    Description string
}

type AdminInfoConfigType struct {
    Name string
    Description string
    Email string
}

type ListenerConfigType struct {
    Host string
    Port int
    Type string
    Tls bool
}

type OperatorConfigType struct {
    Name string
    Matches []string
    Password PasswordConfigType
}

type PasswordConfigType struct {
    Hash string
    Type string
}

func NewConfigurationFile(content []byte) (*ConfigurationFile, os.Error) {
    var c ConfigurationFile
    jsonError := json.Unmarshal(content, &c)
    return &c, jsonError
}
