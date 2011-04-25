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
    "bufio"
    "flag"
    "fmt"
    "os"
    "strings"
    "websocket"
)

var host = flag.String("host", "", "Host to connect to")
var origin = flag.String("origin", "", "Origin to sent")

func Strip(s string) string {
    return strings.Trim(s, " \r\n")
}

func main() {
    flag.Parse()

    if len(*host) == 0 {
        fmt.Print("Error: Needs -host parameter.\n")
        os.Exit(1)
    }

    ws, err := websocket.Dial(*host, "", *origin)

    if err != nil {
        fmt.Printf("Error: %s\n", err.String())
        os.Exit(1)
    }

    rw := bufio.NewReadWriter(bufio.NewReader(ws), bufio.NewWriter(ws))
    stdin := bufio.NewReader(os.Stdin);

    go func() {
        for {
            s, _ := rw.ReadString('\n')
            s = Strip(s)

            if len(s) == 0 {
                continue
            }

            fmt.Printf("%s\n", s)
        }
    }()

    for {
        input, _ := stdin.ReadString('\n');
        rw.WriteString(input)
        rw.Flush()
    }
}
