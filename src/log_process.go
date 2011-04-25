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
    Debug LogLevel = iota
    Warning
    Information
    Silent
)

func (l LogLevel) String() string {
    switch l {
        case Silent: return "Silent"
        case Information: return "Information"
        case Warning: return "Warning"
        case Debug: return "Debug"
    }

    panic("Unhandled Loglevel")
}

type Log struct {
    process *LogProcess
    log_level LogLevel
}

func NewLog(filename string) *Log {
    return &Log{NewLogProcess(filename), Debug}
}

func (log *Log) IsRunning() bool {
    r := make(chan bool)
    go log.process.IsRunning(r)
    return <-r
}

func (log *Log) SetLogLevel(ll LogLevel) {
    log.log_level = ll
}

func (log *Log) Log(ll LogLevel, message string) {
    log.LogF(ll, "%s", message)
}

func (log *Log) LogF(ll LogLevel, format string, arguments ...interface{}) {
    if ll == Silent {
        panic("Using LogLevel 'Silent' does not make sense here")
    }

    if ll >= log.log_level && log.IsRunning() {
        log.process.log_chan <- formatMessage(ll, format, arguments...)
    }
}

func (log *Log) Shutdown() {
    if log.IsRunning() {
        log.process.stop_chan <- true
    } else {
        panic("Trying to Shutdown() a dead LogProcess")
    }
}

type LogProcess struct {
    stop_chan chan bool
    log_chan chan string

    file *os.File
    running bool
    logger *log.Logger
}

func NewLogProcess(filename string) *LogProcess {
    file, error := os.Open(filename, os.O_APPEND | os.O_WRONLY | os.O_CREAT, 0644)

    if error != nil {
        panic("Error: " + error.String())
    }

    logger := log.New(file, "", log.Ldate | log.Ltime)

    stop_chan := make(chan bool)
    log_chan  := make(chan string)

    lp := &LogProcess{stop_chan, log_chan, file, false, logger}
    go func() {
        if lp.running {
            panic("Trying to start an already running LogProcess")
        }

        lp.running = true
        defer lp.Shutdown()

        for {
            select {
                case message := <-lp.log_chan:
                    lp.logger.Print(message)

                case <-lp.stop_chan:
                    return
            }
        }
    }()

    return lp
}

func (process *LogProcess) IsRunning(r chan bool) {
    r <-process.running
}

func (process *LogProcess) Shutdown() {
    defer process.file.Close()
    process.logger.Print(formatMessage(Debug, "Shutting Down Logging Process"))
    process.running = false
}

func formatMessage(ll LogLevel, format string, arguments ...interface{}) string {
    return "(" + ll.String() + "): " + fmt.Sprintf(format, arguments...)
}
