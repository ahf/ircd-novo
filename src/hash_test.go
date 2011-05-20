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

type hashTestCase struct {
    input string
    hash_function string
    expected string
}

var hashTests = []hashTestCase {
    hashTestCase{"", "md5", "d41d8cd98f00b204e9800998ecf8427e"},
    hashTestCase{"", "sha1", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
    hashTestCase{"", "sha256", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
    hashTestCase{"", "sha512", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},

    hashTestCase{"foobar", "MD5", "3858f62230ac3c915f300c664312c63f"},
    hashTestCase{"foobar", "SHA1", "8843d7f92416211de9ebb963ff4ce28125932878"},
    hashTestCase{"foobar", "SHA256", "c3ab8ff13720e8ad9047dd39466b3c8974e592c2fa383d4a3960714caef0c4f2"},
    hashTestCase{"foobar", "SHA512", "0a50261ebd1a390fed2bf326f2673c145582a6342d523204973d0219337f81616a8069b012587cf5635f6925f1b56c360230c19b273500ee013e030601bf2425"},

    hashTestCase{"Færøy", "mD5", "2f9590d271f68585a9bccfaba90e2eb5"},
    hashTestCase{"Goat", "shA1", "0e3f5fc25846b7ffc7fbf7ea4a19522f60c6c683"},
    hashTestCase{"Goat", "sHA256", "ac68bd931568ef75452d7594f2facc95011231b637ca848b2135b6b5aeb66c52"},
    hashTestCase{"Goats", "SHa512", "0ec6fd985dad4e6fb1ccd88e7a5454ae0e190381881dbf49da03efb4acc7dcd3678f12eee2d569c9fa6b1ff2b5d03cb9162d46671dd01201eb4cab7b6fd5899d"},

    hashTestCase{"foobar", "magichashfunction", ""},
    hashTestCase{"XXXxxxXXX", "", ""},
}

func TestHashFunctionFromString(t *testing.T) {

}

func TestHash(t *testing.T) {
    for i := range hashTests {
        test := hashTests[i]
        result := Hash(test.hash_function, test.input)

        if result != test.expected {
            t.Errorf("Hash('%s', '%s') = '%s', want '%s'.", test.hash_function, test.input, result, test.expected)
        }
    }
}
