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
    "fmt"
    "log"
    "os"
)

type LogLevel int

const (
    Silent LogLevel = iota
    Warning
    Debug
)

func (l LogLevel) String() string {
    switch l {
        case Silent: return "Silent"
        case Warning: return "Warning"
        case Debug: return "Debug"
    }

    panic("Unhandled Loglevel")
}

func CreateLogMessage(ll LogLevel, format string, arguments ...interface{}) string {
    return "(" + ll.String() + "): " + fmt.Sprintf(format, arguments...)
}

type LogProcess struct {
    logger *log.Logger
    file *os.File
    command_channel chan logProcessCommand
    logging_channel chan string
    valid bool
    log_level LogLevel
}

func NewLogProcess(filename string) *LogProcess {
    file, error := os.Open(filename, os.O_APPEND | os.O_WRONLY | os.O_CREAT, 0644)

    if error != nil {
        panic("Error: " + error.String())
    }

    logger := log.New(file, nil, "", log.Ldate | log.Ltime)
    command_channel := make(chan logProcessCommand)
    logging_channel := make(chan string)

    lp := LogProcess{logger, file, command_channel, logging_channel, true, Silent}

    go func() {
        for {
            select {
                case message := <-lp.logging_channel:
                    lp.logger.Log(message)
                case command := <-lp.command_channel:
                    lp.logger.Log(CreateLogMessage(Debug, "Received '%s' Command Message", command.String()))

                    switch command {
                        case closeCommand:
                            return
                    }
            }
        }
    }()

    return &lp
}

func (s *LogProcess) LogF(ll LogLevel, format string, arguments ...interface{}) {
    if ll == Silent {
        panic("Using LogLevel 'Silent' does not make sense here")
    }

    if s.valid && ll >= s.log_level {
        s.logging_channel <- CreateLogMessage(ll, format, arguments...)
    }
}

func (s *LogProcess) Log(ll LogLevel, message string) {
    s.LogF(ll, "%s", message)
}

func (s *LogProcess) SetLogLevel(ll LogLevel) {
    s.log_level = ll
}

func (s *LogProcess) Close() {
    if s.valid {
        s.command_channel <- closeCommand
        s.valid = false
        s.file.Close()
    }
}

type logProcessCommand int

const (
    closeCommand logProcessCommand = iota
)

func (l logProcessCommand) String() string {
    switch l {
        case closeCommand: return "Close"
    }

    panic("Unhandled logProcessCommand")
}
