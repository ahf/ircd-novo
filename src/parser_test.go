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
    "testing"
)

func stringSpliceCompare(a, b []string) bool {
    if len(a) != len(b) {
        return false
    }

    for i := range a {
        if a[i] != b[i] {
            return false
        }
    }

    return true
}

type parserTestCase struct {
    input string
    expected_prefix string
    expected_command string
    expected_arguments []string
}

var parserTestCases = []parserTestCase {
    parserTestCase{"", "", "", []string{}},
    parserTestCase{"NICK ahf", "", "NICK", []string{"ahf"}},
    parserTestCase{"USER ahf ahf irc.0x90.dk :Alexander Færøy", "", "USER", []string{"ahf", "ahf", "irc.0x90.dk", "Alexander Færøy"}},
    parserTestCase{"NOTICE AUTH :*** Looking up your hostname...", "", "NOTICE", []string{"AUTH", "*** Looking up your hostname..."}},
    parserTestCase{":irc.0x90.dk 004 ahf_ irc.0x90.dk ircd-ratbox-3.0.6 oiwszcrkfydnxbauglZCD biklmnopstveIrS bkloveI", "irc.0x90.dk", "004", []string{"ahf_", "irc.0x90.dk", "ircd-ratbox-3.0.6", "oiwszcrkfydnxbauglZCD", "biklmnopstveIrS", "bkloveI"}},
}

func TestParser(t *testing.T) {
    for i := range parserTestCases {
        test := parserTestCases[i]
        result := Parse(test.input)

        if result.Prefix() != test.expected_prefix {
            t.Errorf("Parse('%s') Prefix = '%s', want '%s'.", test.input, result.Prefix(), test.expected_prefix)
        }

        if result.Command() != test.expected_command {
            t.Errorf("Parse('%s') Command = '%s', want '%s'.", test.input, result.Command(), test.expected_command)
        }

        if ! stringSpliceCompare(result.Arguments(), test.expected_arguments) {
            t.Errorf("Parse('%s') Arguments = '%s', want '%s'.", test.input, Join(result.Arguments()), Join(test.expected_arguments))
        }
    }
}
