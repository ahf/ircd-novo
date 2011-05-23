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

type testCase struct {
    s string
    r bool
}

var nickNameTests = []testCase {
    testCase{"", false},
    testCase{"[[(0)]]", true},
    testCase{"abc", true},
    testCase{"ahf", true},
    testCase{"foobar", true},
    testCase{"æøå", false},
}

var userNameTests = []testCase {
    testCase{"", false},
    testCase{"abc", true},
    testCase{"ahf", true},
    testCase{"æøå", false},
    testCase{"{}", false},
    testCase{"|||", false},
}

var hostNameTests = []testCase {
    testCase{"", false},
    testCase{"foo.bar.baz", true},
    testCase{"foobar", true},
    testCase{"irssi.org", true},
    testCase{"horus.0x90.dk", true},
    testCase{"really.annoyingly.long.host.name.that.idiots.use.on.irc", true}, // :-(
}

type collapseTestCase struct {
    input string
    output string
}

var collapseTests = []collapseTestCase {
    collapseTestCase{"", ""},
    collapseTestCase{"***", "*"},
    collapseTestCase{"* * *", "* * *"},
    collapseTestCase{"*foobar**", "*foobar*"},
}

func TestIsValidNickname(t *testing.T) {
    for i := range nickNameTests {
        test := nickNameTests[i]

        if IsValidNickname(test.s) != test.r {
            t.Errorf("Validation of nickname '%s' failed.", test.s)
        }
    }
}

func TestIsValidUsername(t *testing.T) {
    for i := range userNameTests {
        test := userNameTests[i]

        if IsValidUsername(test.s) != test.r {
            t.Errorf("Validation of username '%s' failed.", test.s)
        }
    }
}

func TestIsValidHostname(t *testing.T) {
    for i := range hostNameTests {
        test := hostNameTests[i]

        if IsValidUsername(test.s) != test.r {
            t.Errorf("Validation of hostname '%s' failed.", test.s)
        }
    }
}

func TestCollapse(t *testing.T) {
    for i := range collapseTests {
        test := collapseTests[i]
        result := Collapse(test.input)

        if result != test.output {
            t.Errorf("Collapse('%s') = '%s', want '%s'.", test.input, result, test.output)
        }
    }
}
